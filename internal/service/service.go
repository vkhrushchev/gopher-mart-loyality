package service

import (
	"context"
	"errors"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

var (
	ErrWithdrawNoFundsOnBalance = errors.New("service_withdraw: no found on balance")
)

type IUserService interface {
	RegisterUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, password string) (string, error)
	GetBalance(ctx context.Context) (dto.UserBalanceDomain, error)
}

type IOrderService interface {
	PutOrder(ctx context.Context, orderNumber string) (bool, error)
	GetOrders(ctx context.Context) ([]dto.OrderDomain, error)
}

type IWithDrawService interface {
	MakeWithdraw(ctx context.Context, orderNumber string, sum float64) error
	GetUserWithdraws(ctx context.Context) ([]dto.UserWithdrawDomain, error)
}
