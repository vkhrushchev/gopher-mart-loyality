package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
	mock_storage "github.com/vkhrushchev/gopher-mart-loyality/internal/storage/mock"
)

func TestWithdrawalService_DoWithdrawal(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	orderStorageMock := mock_storage.NewMockIOrderStorage(mockController)
	userStorageMock := mock_storage.NewMockIUserStorage(mockController)
	withdrawalStorageMock := mock_storage.NewMockIWithdrawalStorage(mockController)
	withdrawalService := NewWithdrawalService(orderStorageMock, userStorageMock, withdrawalStorageMock)

	type args struct {
		ctx         context.Context
		orderNumber string
		sum         float64
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(
			orderStorageMock *mock_storage.MockIOrderStorage,
			userStorageMock *mock_storage.MockIUserStorage,
			withdrawalStorageMock *mock_storage.MockIWithdrawalStorage)
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: "0123456789",
				sum:         10.5,
			},
			prepareMocks: func(
				orderStorageMock *mock_storage.MockIOrderStorage,
				userStorageMock *mock_storage.MockIUserStorage,
				withdrawalStorageMock *mock_storage.MockIWithdrawalStorage) {
				orderStorageMock.EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(
						&dto.OrderEntity{},
						nil)
				userStorageMock.EXPECT().
					GetUserBalanceByLogin(gomock.Any(), gomock.Any()).
					Return(
						&dto.UserBalanceEntity{
							TotalSum:           100.5,
							TotalWithdrawalSum: 0.0,
						},
						nil)
				withdrawalStorageMock.EXPECT().
					SaveBalanceWithdrawal(gomock.Any(), gomock.Any()).
					Return(&dto.BalanceWithdrawalEntity{}, nil)
			},
			expectedErr: nil,
		},
		{
			name: "wrong order number",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: "0123456789",
				sum:         10.5,
			},
			prepareMocks: func(
				orderStorageMock *mock_storage.MockIOrderStorage,
				userStorageMock *mock_storage.MockIUserStorage,
				withdrawalStorageMock *mock_storage.MockIWithdrawalStorage) {
				orderStorageMock.EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrEntityNotFound)
			},
			expectedErr: ErrOrderWrongNumber,
		},
		{
			name: "no funds on balance",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: "0123456789",
				sum:         10.5,
			},
			prepareMocks: func(
				orderStorageMock *mock_storage.MockIOrderStorage,
				userStorageMock *mock_storage.MockIUserStorage,
				withdrawalStorageMock *mock_storage.MockIWithdrawalStorage) {
				orderStorageMock.EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(
						&dto.OrderEntity{},
						nil)
				userStorageMock.EXPECT().
					GetUserBalanceByLogin(gomock.Any(), gomock.Any()).
					Return(
						&dto.UserBalanceEntity{
							TotalSum:           0.0,
							TotalWithdrawalSum: 0.0,
						},
						nil)
			},
			expectedErr: ErrNoFundsOnBalance,
		},
		{
			name: "storage no funds on balance",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: "0123456789",
				sum:         10.5,
			},
			prepareMocks: func(
				orderStorageMock *mock_storage.MockIOrderStorage,
				userStorageMock *mock_storage.MockIUserStorage,
				withdrawalStorageMock *mock_storage.MockIWithdrawalStorage) {
				orderStorageMock.EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(
						&dto.OrderEntity{},
						nil)
				userStorageMock.EXPECT().
					GetUserBalanceByLogin(gomock.Any(), gomock.Any()).
					Return(
						&dto.UserBalanceEntity{
							TotalSum:           100.5,
							TotalWithdrawalSum: 0.0,
						},
						nil)
				withdrawalStorageMock.EXPECT().
					SaveBalanceWithdrawal(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrNoFundsOnBalance)
			},
			expectedErr: ErrNoFundsOnBalance,
		},
	}
	for _, tt := range tests {
		tt.prepareMocks(orderStorageMock, userStorageMock, withdrawalStorageMock)

		t.Run(tt.name, func(t *testing.T) {
			err := withdrawalService.DoWithdrawal(tt.args.ctx, tt.args.orderNumber, tt.args.sum)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}
		})
	}
}

func TestWithdrawalService_GetUserWithdrawals(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	orderStorage := mock_storage.NewMockIOrderStorage(mockController)
	userStorageMock := mock_storage.NewMockIUserStorage(mockController)
	withdrawalStorageMock := mock_storage.NewMockIWithdrawalStorage(mockController)
	withdrawalService := NewWithdrawalService(orderStorage, userStorageMock, withdrawalStorageMock)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(
			orderStorageMock *mock_storage.MockIOrderStorage,
			userStorageMock *mock_storage.MockIUserStorage,
			withdrawalStorageMock *mock_storage.MockIWithdrawalStorage)
		expectedResult []dto.OrderWithdrawalDomain
		expectedErr    error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
			},
			prepareMocks: func(
				orderStorageMock *mock_storage.MockIOrderStorage,
				userStorageMock *mock_storage.MockIUserStorage,
				withdrawalStorageMock *mock_storage.MockIWithdrawalStorage) {
				withdrawalStorageMock.EXPECT().
					GetBalanceWithdrawalsByUserLogin(gomock.Any(), gomock.Any()).
					Return(
						[]dto.BalanceWithdrawalEntity{
							{
								OrderNumber:   "0000111122223333",
								WithdrawalSum: 10.5,
								ProcessedAt:   time.Date(2024, time.September, 13, 15, 32, 0, 0, time.UTC),
							},
						},
						nil,
					)
			},
			expectedResult: []dto.OrderWithdrawalDomain{
				{
					OrderNumber:   "0000111122223333",
					WithdrawalSum: 10.5,
					ProcessedAt:   time.Date(2024, time.September, 13, 15, 32, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(orderStorage, userStorageMock, withdrawalStorageMock)

			result, err := withdrawalService.GetUserWithdrawals(tt.args.ctx)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
