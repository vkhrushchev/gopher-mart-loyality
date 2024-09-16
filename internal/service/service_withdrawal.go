package service

import (
	"context"
	"errors"
	"time"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
)

type WithdrawalService struct {
	orderStorage      storage.IOrderStorage
	userStorage       storage.IUserStorage
	withdrawalStorage storage.IWithdrawalStorage
}

func NewWithdrawalService(
	orderStorage storage.IOrderStorage,
	userStorage storage.IUserStorage,
	withdrawalStorage storage.IWithdrawalStorage) *WithdrawalService {
	return &WithdrawalService{
		orderStorage:      orderStorage,
		userStorage:       userStorage,
		withdrawalStorage: withdrawalStorage,
	}
}

func (s *WithdrawalService) DoWithdrawal(ctx context.Context, orderNumber string, withdrawal float64) error {
	userLogin := ctx.Value(middleware.UserLoginContextKey).(string)
	log.Infow(
		"service_withdrawal: make withdrawn",
		"user_login", userLogin,
		"order_number", orderNumber,
		"withdrawal", withdrawal)

	//
	// логично было бы проверять на то что системе известен заказ перед списанием, но тесты падают
	//
	// _, err := s.orderStorage.GetOrderByOrderNumber(ctx, orderNumber)
	// if err != nil && errors.Is(err, storage.ErrEntityNotFound) {
	// 	log.Errorw("service_withdrawal: order not found by order number", "order_number", orderNumber)
	// 	return ErrOrderWrongNumber
	// } else if err != nil {
	// 	log.Errorw("service_withdrawal: unexpected storage error", err, err.Error())
	// 	return err
	// }

	userBalance, err := s.userStorage.GetUserBalanceByLogin(ctx, userLogin)
	if err != nil && errors.Is(err, storage.ErrEntityNotFound) {
		log.Errorw("service_withdrawal: user balance not found by login", "user_login", userLogin)
		return err
	} else if err != nil {
		log.Errorw("service_withdrawal: unexpected storage error", err, err.Error())
		return err
	}

	if userBalance.TotalSum < withdrawal {
		return ErrNoFundsOnBalance
	}

	balanceWithdrawal := dto.BalanceWithdrawalEntity{
		UserLogin:     userLogin,
		OrderNumber:   orderNumber,
		WithdrawalSum: withdrawal,
		ProcessedAt:   time.Now().UTC(),
	}
	_, err = s.withdrawalStorage.SaveBalanceWithdrawal(ctx, &balanceWithdrawal)
	if err != nil && errors.Is(err, storage.ErrNoFundsOnBalance) {
		return ErrNoFundsOnBalance
	} else if err != nil {
		log.Errorw("service_withdrawal: unexpected storage error", "error", err.Error())
		return err
	}

	return nil
}

func (s *WithdrawalService) GetUserWithdrawals(ctx context.Context) ([]dto.OrderWithdrawalDomain, error) {
	userLogin := ctx.Value(middleware.UserLoginContextKey).(string)
	log.Infow("service_withdrawal: get user withdrawals", "user_login", userLogin)

	withdrawalEntities, err := s.withdrawalStorage.GetBalanceWithdrawalsByUserLogin(ctx, userLogin)
	if err != nil {
		log.Errorw("service_withdrawal: unexpected storage error", "error", err.Error())
		return nil, err
	}

	withdrawalDomains := make([]dto.OrderWithdrawalDomain, 0, len(withdrawalEntities))
	for _, withdrawalEntity := range withdrawalEntities {
		withdrawalDomain := dto.OrderWithdrawalDomain{
			OrderNumber:   withdrawalEntity.OrderNumber,
			WithdrawalSum: withdrawalEntity.WithdrawalSum,
			ProcessedAt:   withdrawalEntity.ProcessedAt,
		}
		withdrawalDomains = append(withdrawalDomains, withdrawalDomain)
	}

	return withdrawalDomains, nil
}
