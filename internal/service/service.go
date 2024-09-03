package service

import (
	"context"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

type IUserService interface {
	RegisterUser(ctx context.Context, username string, password string) error
	LoginUser(ctx context.Context, username string, password string) (string, error)
	GetBalance(ctx context.Context) (dto.UserBalanceDomain, error)
}
