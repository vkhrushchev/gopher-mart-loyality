package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

var ErrAccrualInternalServerError = fmt.Errorf("service_accrual: internal server error")
var ErrAccrualUnknownOrder = fmt.Errorf("service_accrual: unknown order")
var ErrAccrualRateLimitExceed = fmt.Errorf("servcice_accrual: rate limit exceed")

type AccrualService struct {
	accrualURL string
	client     *http.Client
}

func NewAccrualService(accrualURL string) *AccrualService {
	client := &http.Client{}

	return &AccrualService{
		accrualURL: accrualURL,
		client:     client,
	}
}

func (s *AccrualService) GetAccrualInfo(ctx context.Context, orderNumber string) (*dto.AccuralInfoDomain, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, s.accrualURL+"/api/orders/"+orderNumber, nil)
	if err != nil {
		log.Errorw("service_accrual: error when build request", "error", err.Error())
		return nil, err
	}

	res, err := s.client.Do(r)
	if err != nil {
		log.Errorw("service_accrual: error when do request", "error", err.Error())
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Errorw("accrual_service: error when close response body")
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil, ErrAccrualUnknownOrder
	}

	if res.StatusCode == http.StatusTooManyRequests {
		return nil, ErrAccrualRateLimitExceed
	}

	if res.StatusCode == http.StatusInternalServerError {
		return nil, ErrAccrualInternalServerError
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Errorw("service_accrual: unexpected error when read body", "error", err.Error())
		return nil, err
	}

	var accrualOrderResponse dto.AccuralOrderResponse
	if err := json.Unmarshal(bodyBytes, &accrualOrderResponse); err != nil {
		log.Errorw("service_accrual: unexpected error when parse body", "error", err.Error())
		return nil, err
	}

	accuralInfoDomain := dto.AccuralInfoDomain{
		OrderNumber: accrualOrderResponse.Order,
		OrderStatus: accrualOrderResponse.Status,
		Accrual:     accrualOrderResponse.Accrual,
	}

	return &accuralInfoDomain, nil
}
