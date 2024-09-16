package service

import (
	"context"
	"errors"
	"time"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

type AccrualPullerService struct {
	accrualService                IAccrualService
	orderService                  IOrderService
	startProcessingOrderCh        chan string
	pullAccrualServiceCh          chan string
	processAccrualServiceResultCh chan pullAccrualServiceResult
}

type pullAccrualServiceResult struct {
	orderNumber   string
	accrualStatus dto.AccrualStatus
	accrual       float64
}

func NewAccrualPullerService(accrualService IAccrualService, orderService IOrderService) *AccrualPullerService {
	return &AccrualPullerService{
		accrualService:                accrualService,
		orderService:                  orderService,
		startProcessingOrderCh:        make(chan string, 100),
		pullAccrualServiceCh:          make(chan string, 100),
		processAccrualServiceResultCh: make(chan pullAccrualServiceResult, 100),
	}
}

func (s *AccrualPullerService) Start() {
	s.startProcessAccrualServiceResult()
	s.startPullAccrualService()
	s.startProcessingOrder()
}

func (s *AccrualPullerService) startProcessingOrder() {
	go func() {
		for {
			orderNumber := <-s.startProcessingOrderCh
			log.Infow("service_accrual_puller: start processing order", "order_number", orderNumber)

			err := s.orderService.UpdateOrderStatus(context.Background(), orderNumber, dto.OrderStatusProcessing)
			if err != nil {
				log.Errorw("service_accrual_puller: unexpected error when update order status")
			}

			s.pullAccrualServiceCh <- orderNumber
		}
	}()
}

func (s *AccrualPullerService) startPullAccrualService() {
	go func() {
		for {
			orderNumber := <-s.pullAccrualServiceCh
			log.Infow("service_accrual_puller: start accruall pulling", "order_number", orderNumber)

			accrualInfo, err := s.accrualService.GetAccrualInfo(context.Background(), orderNumber)
			if err != nil {
				if errors.Is(err, ErrAccrualUnknownOrder) {
					// предполагал что если заказ неизвестен сервису расчетов бонусов то это финальная ситуация
					// log.Errorw(
					// 	"service_accrual_puller: order number unknown in accrual service. update status to \"INVALID\" and skip",
					// 	"order_number", orderNumber)

					// if err := s.orderService.UpdateOrderStatus(context.Background(), orderNumber, dto.OrderStatusInvalid); err != nil {
					// 	log.Errorw("service_accrual_puller: unexpected error when update order status")
					// }

					// похоже заказ в системе начислений может появиться позже, переотправим c паузой в 1с
					log.Errorw(
						"service_accrual_puller: order number unknown in accrual service. retry after 1s pause",
						"order_number", orderNumber)
					go func() {
						time.Sleep(1 * time.Second)
						s.pullAccrualServiceCh <- orderNumber
					}()

					continue
				}

				if errors.Is(err, ErrAccrualInternalServerError) {
					log.Errorw("service_accrual_puller: rate limit exceed. pause for 1m and retry")

					s.pullAccrualServiceCh <- orderNumber
					time.Sleep(1 * time.Minute)

					continue
				}

				if errors.Is(err, ErrAccrualInternalServerError) {
					log.Errorw("service_accrual_puller: accrual internal server error. pause for 100ms and retry")

					s.pullAccrualServiceCh <- orderNumber
					time.Sleep(100 * time.Millisecond)

					continue
				}

				log.Errorw("service_accrual_puller: unexpected accrual error. skip", "error", err.Error())
			}

			s.processAccrualServiceResultCh <- pullAccrualServiceResult{
				orderNumber:   accrualInfo.OrderNumber,
				accrualStatus: accrualInfo.OrderStatus,
				accrual:       accrualInfo.Accrual,
			}
		}
	}()
}

func (s *AccrualPullerService) startProcessAccrualServiceResult() {
	go func() {
		for {
			pullAccrualServiceResult := <-s.processAccrualServiceResultCh

			orderNumber := pullAccrualServiceResult.orderNumber
			accrualStatus := pullAccrualServiceResult.accrualStatus
			accrual := pullAccrualServiceResult.accrual
			log.Infow("service_accrual_puller: start processing pulling result", "order_number", orderNumber)

			switch accrualStatus {
			case dto.AccrualRegistredStatus:
			case dto.AccrualProcessingStatus:
				go func() {
					log.Infow("service_accrual_puller: retry accrual pulling", "order_number", orderNumber)

					time.Sleep(1 * time.Second)
					s.pullAccrualServiceCh <- pullAccrualServiceResult.orderNumber
				}()
			case dto.AccrualInvalidStatus:
				log.Infow("service_accrual_puller: update order", "order_number", orderNumber)
				err := s.orderService.UpdateOrderStatus(context.Background(), orderNumber, dto.OrderStatusInvalid)
				if err != nil {
					log.Errorw("service_accrual_puller: unexpected error when update order status")
				}
			case dto.AccrualProcessedStatus:
				log.Infow("service_accrual_puller: update order", "order_number", orderNumber)
				err := s.orderService.UpdateOrderStatusAndAccrual(context.Background(), orderNumber, dto.OrderStatusProcessed, accrual)
				if err != nil {
					log.Errorw("service_accrual_puller: unexpected error when update order status and accrual")
				}
			}
		}
	}()
}

func (s *AccrualPullerService) AddGetAccrualInfoTask(ctx context.Context, orderNumber string) {
	s.startProcessingOrderCh <- orderNumber
}
