package storage

import (
	"context"
	"errors"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

var (
	ErrUnexpextedDBError = errors.New("storage: unexpected DB error")
	ErrEntityExists      = errors.New("storage: entity exists")
	ErrEntityNotFound    = errors.New("storage: no entity found")
)

type IUserStorage interface {
	SaveUser(ctx context.Context, user *dto.UserEntity) (*dto.UserEntity, error)
	GetUserByLoginAndPasswordHash(ctx context.Context, login string, passwordHash string) (*dto.UserEntity, error)
}

type IOrderStorage interface {
	SaveOrder(ctx context.Context, order *dto.OrderEntity) (*dto.OrderEntity, error)
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (*dto.OrderEntity, error)
	GetOrdersByUserLogin(ctx context.Context, userLogin string) ([]dto.OrderEntity, error)
}
