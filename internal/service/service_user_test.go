package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
	mock_storage "github.com/vkhrushchev/gopher-mart-loyality/internal/storage/mock"
)

var (
	salt         string = "salt"
	jwtSecretKey string = "jwtSecretKey"
)

func TestUserService_RegisterUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userStorageMock := mock_storage.NewMockIUserStorage(mockController)
	userService := NewUserService(userStorageMock, salt, jwtSecretKey)

	type args struct {
		ctx      context.Context
		username string
		password string
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(userStorageMock *mock_storage.MockIUserStorage)
		expectedErr  error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			prepareMocks: func(userStorageMock *mock_storage.MockIUserStorage) {
				userStorageMock.
					EXPECT().
					SaveUser(gomock.Any(), gomock.Any()).
					Return(&dto.UserEntity{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "user exists",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			prepareMocks: func(userStorageMock *mock_storage.MockIUserStorage) {
				userStorageMock.
					EXPECT().
					SaveUser(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrEntityExists)
			},
			expectedErr: ErrUserExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(userStorageMock)
			err := userService.RegisterUser(tt.args.ctx, tt.args.username, tt.args.password)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userStorageMock := mock_storage.NewMockIUserStorage(mockController)
	userService := NewUserService(userStorageMock, salt, jwtSecretKey)

	type args struct {
		ctx      context.Context
		username string
		password string
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(userStorageMock *mock_storage.MockIUserStorage)
		expectedErr  error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			prepareMocks: func(userStorageMock *mock_storage.MockIUserStorage) {
				userStorageMock.
					EXPECT().
					GetUserByLoginAndPasswordHash(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&dto.UserEntity{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			prepareMocks: func(userStorageMock *mock_storage.MockIUserStorage) {
				userStorageMock.
					EXPECT().
					GetUserByLoginAndPasswordHash(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrEntityNotFound)
			},
			expectedErr: ErrWrongLoginOrPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(userStorageMock)

			token, err := userService.LoginUser(tt.args.ctx, tt.args.username, tt.args.password)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			} else if errors.Is(err, tt.expectedErr) {
				return
			}

			assert.NotNil(t, token)
		})
	}
}

func TestUserService_GetBalance(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userStorageMock := mock_storage.NewMockIUserStorage(mockController)
	userService := NewUserService(userStorageMock, salt, jwtSecretKey)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		prepareMocks   func(userStorageMock *mock_storage.MockIUserStorage)
		expectedResult dto.UserBalanceDomain
		expectedErr    error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
			},
			prepareMocks: func(userStorageMock *mock_storage.MockIUserStorage) {
				userStorageMock.EXPECT().
					GetUserBalanceByLogin(gomock.Any(), gomock.Any()).
					Return(
						&dto.UserBalanceEntity{
							TotalSum:           10.5,
							TotalWithdrawalSum: 0.5,
						},
						nil)
			},
			expectedResult: dto.UserBalanceDomain{
				Current:    10.5,
				Withdrawal: 0.5,
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(userStorageMock)

			result, err := userService.GetBalance(tt.args.ctx)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			} else if errors.Is(err, tt.expectedErr) {
				return
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
