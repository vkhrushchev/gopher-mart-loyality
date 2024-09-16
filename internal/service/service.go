package service

import (
	"context"
	"errors"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

var (
	ErrNoFundsOnBalance = errors.New("service_withdraw: no found on balance")
)

type IUserService interface {
	RegisterUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, password string) (string, error)
	GetBalance(ctx context.Context) (dto.UserBalanceDomain, error)
}

type IOrderService interface {
	PutOrder(ctx context.Context, orderNumber string) (bool, error)
	UpdateOrderStatus(ctx context.Context, orderNumber string, orderStatus dto.OrderStatus) error
	UpdateOrderStatusAndAccrual(ctx context.Context, orderNumber string, orderStatus dto.OrderStatus, accrual float64) error
	GetOrders(ctx context.Context) ([]dto.OrderDomain, error)
}

type IWithdrawalService interface {
	DoWithdrawal(ctx context.Context, orderNumber string, sum float64) error
	GetUserWithdrawals(ctx context.Context) ([]dto.OrderWithdrawalDomain, error)
}

type IAccrualService interface {
	GetAccrualInfo(ctx context.Context, orderNumber string) (*dto.AccuralInfoDomain, error)
}

type IAccrualPullerService interface {
	Start()
	AddGetAccrualInfoTask(ctx context.Context, orderNumber string)
}
