package service

import (
	"context"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

type WithdrawService struct {
}

func NewWithdrawService() *WithdrawService {
	return &WithdrawService{}
}

func (s *WithdrawService) MakeWithdraw(ctx context.Context, orderNumber string, sum float64) error {
	log.Infow(
		"make withdrawn.",
		"number", orderNumber,
		"sum", sum)

	return nil
}

func (s *WithdrawService) GetUserWithdraws(ctx context.Context) ([]dto.UserWithdrawDomain, error) {
	log.Infow("get user withdraws.")

	return make([]dto.UserWithdrawDomain, 0), nil
}
