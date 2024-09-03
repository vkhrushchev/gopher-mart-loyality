package service

import (
	"context"
	"fmt"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

var ErrUnknownOrder = fmt.Errorf("service_accrual: unknown order")
var ErrRateLimitExceed = fmt.Errorf("servcice_accrual: rate limit exceed")

type AccrualService struct {
}

func NewAccrualService() *AccrualService {
	return &AccrualService{}
}

func (s *AccrualService) GetAccrualInfo(ctx context.Context, orderNumber string) (dto.AccuralInfoDomain, error) {
	return dto.AccuralInfoDomain{}, nil
}
