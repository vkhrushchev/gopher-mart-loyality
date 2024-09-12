package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
)

var (
	ErrOrderUploadedByAnotherUser = fmt.Errorf("order already uploaded by another user")
	ErrOrderWrongNumber           = fmt.Errorf("order wrong number")
)

type OrderService struct {
	orderStorage storage.IOrderStorage
}

func NewOrderService(orderStorage storage.IOrderStorage) *OrderService {
	return &OrderService{
		orderStorage: orderStorage,
	}
}

func (s *OrderService) PutOrder(ctx context.Context, orderNumber string) (bool, error) {
	userLogin := ctx.Value(middleware.UserLoginContextKey).(string)
	log.Infow(
		"service_order: put order",
		"order_number", orderNumber,
		"user_login", userLogin,
	)

	if err := goluhn.Validate(orderNumber); err != nil {
		return false, ErrOrderWrongNumber
	}

	orderEntity := &dto.OrderEntity{
		UserLogin:  userLogin,
		Number:     orderNumber,
		Status:     dto.OrderStatusNew,
		UploadedAt: time.Now().UTC(),
	}
	_, err := s.orderStorage.SaveOrder(ctx, orderEntity)
	if err != nil && errors.Is(err, storage.ErrEntityExists) {
		existedOrderEntity, err := s.orderStorage.GetOrderByOrderNumber(ctx, orderNumber)
		if err != nil {
			log.Errorw(
				"service_order: unexpected storage error",
				"error", err.Error(),
			)

			return false, nil
		}

		if existedOrderEntity.UserLogin != userLogin {
			return false, ErrOrderUploadedByAnotherUser
		}

		return true, nil
	}

	return false, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]dto.OrderDomain, error) {
	log.Infow("service: get order list for user")

	userLogin := ctx.Value(middleware.UserLoginContextKey).(string)
	orderEntities, err := s.orderStorage.GetOrdersByUserLogin(ctx, userLogin)
	if err != nil {
		log.Errorw(
			"service_order: unexpected storage error",
			"error", err.Error(),
		)

		return nil, err
	}

	orderDomains := make([]dto.OrderDomain, 0, len(orderEntities))
	for _, orderEntity := range orderEntities {
		orderDomain := dto.OrderDomain{
			Number:     orderEntity.Number,
			Status:     orderEntity.Status,
			UploadedAt: orderEntity.UploadedAt,
		}

		if orderEntity.Accrual.Valid {
			orderDomain.Accrual = &orderEntity.Accrual.Float64
		}

		orderDomains = append(orderDomains, orderDomain)
	}

	return orderDomains, nil
}
