package service

import (
	"context"
	"errors"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

var (
	ErrWrongLoginOrPassword = errors.New("unknows user or password")
	ErrUserExists           = errors.New("user exists")
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) RegisterUser(ctx context.Context, username string, password string) error {
	log.Infow(
		"service: register user.",
		"username", username)

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, username string, password string) (string, error) {
	log.Infow(
		"service: login user.",
		"username", username)

	return "", nil
}

func (s *UserService) GetBalance(ctx context.Context) (dto.UserBalanceDomain, error) {
	log.Infow(
		"service: get user balance.")

	return dto.UserBalanceDomain{}, nil
}
