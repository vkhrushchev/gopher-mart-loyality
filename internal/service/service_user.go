package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
)

var (
	ErrWrongLoginOrPassword = errors.New("unknows user or password")
	ErrUserExists           = errors.New("user exists")
)

type UserService struct {
	userStorage  storage.IUserStorage
	salt         string
	jwtSecretKey string
}

func NewUserService(userStorage storage.IUserStorage, salt string, jwtSecretKey string) *UserService {
	return &UserService{
		userStorage:  userStorage,
		salt:         salt,
		jwtSecretKey: jwtSecretKey,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, username string, password string) error {
	log.Infow(
		"service: register user.",
		"username", username)

	passwordHashBytes := md5.Sum([]byte(password + s.salt))
	passwordHash := hex.EncodeToString(passwordHashBytes[:])

	userEntity, err := s.userStorage.SaveUser(
		ctx,
		&dto.UserEntity{
			Login:        username,
			PasswordHash: passwordHash,
			Salt:         s.salt,
		})
	if err != nil && errors.Is(err, storage.ErrEntityExists) {
		return ErrUserExists
	} else if err != nil {
		return err
	}

	log.Debugw(
		"service_user: user saved",
		"id", userEntity.Id,
		"username", userEntity.Login,
		"password_hash", userEntity.PasswordHash,
	)

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, username string, password string) (string, error) {
	log.Infow(
		"service: login user.",
		"username", username)

	passwordHashBytes := md5.Sum([]byte(password + s.salt))
	passwordHash := hex.EncodeToString(passwordHashBytes[:])

	_, err := s.userStorage.GetUserByLoginAndPasswordHash(ctx, username, passwordHash)
	if err != nil && errors.Is(err, storage.ErrEntityNotFound) {
		return "", ErrWrongLoginOrPassword
	} else if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "gophermart",
			"sub": username,
		},
	)

	tokenStr, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		log.Errorw("service_user: unexpected error when generage jwt", "error", err.Error())
		return "", err
	}

	return tokenStr, nil
}

func (s *UserService) GetBalance(ctx context.Context) (dto.UserBalanceDomain, error) {
	log.Infow(
		"service: get user balance.")

	return dto.UserBalanceDomain{}, nil
}
