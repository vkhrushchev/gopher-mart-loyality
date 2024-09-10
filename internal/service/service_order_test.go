package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
	mock_storage "github.com/vkhrushchev/gopher-mart-loyality/internal/storage/mock"
)

var (
	testUserLogin string  = "test_user"
	testAccural   float64 = 10.5
)

func TestOrderService_PutOrder(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	orderStorageMock := mock_storage.NewMockIOrderStorage(mockController)
	orderService := NewOrderService(orderStorageMock)

	type args struct {
		ctx         context.Context
		orderNumber string
	}

	tests := []struct {
		name         string
		args         args
		prepareMocks func(orderStorageMock *mock_storage.MockIOrderStorage)
		orderExists  bool
		expectedErr  error
	}{
		{
			name: "sucess",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: goluhn.Generate(16),
			},
			prepareMocks: func(orderStorageMock *mock_storage.MockIOrderStorage) {
				orderStorageMock.
					EXPECT().
					SaveOrder(gomock.Any(), gomock.Any()).
					Return(&dto.OrderEntity{}, nil)
			},
			orderExists: false,
			expectedErr: nil,
		},
		{
			name: "order wrong number",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: "0000000000000001",
			},
			prepareMocks: func(orderStorageMock *mock_storage.MockIOrderStorage) {
			},
			orderExists: false,
			expectedErr: ErrOrderWrongNumber,
		},
		{
			name: "order exists",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: goluhn.Generate(16),
			},
			prepareMocks: func(orderStorageMock *mock_storage.MockIOrderStorage) {
				orderStorageMock.
					EXPECT().
					SaveOrder(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrEntityExists)

				orderStorageMock.
					EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(
						&dto.OrderEntity{
							UserLogin: "test_user",
						},
						nil,
					)
			},
			orderExists: true,
			expectedErr: nil,
		},
		{
			name: "order uploaded by another user",
			args: args{
				ctx:         context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
				orderNumber: goluhn.Generate(16),
			},
			prepareMocks: func(orderStorageMock *mock_storage.MockIOrderStorage) {
				orderStorageMock.
					EXPECT().
					SaveOrder(gomock.Any(), gomock.Any()).
					Return(nil, storage.ErrEntityExists)

				orderStorageMock.
					EXPECT().
					GetOrderByOrderNumber(gomock.Any(), gomock.Any()).
					Return(
						&dto.OrderEntity{
							UserLogin: "not_test_user",
						},
						nil,
					)
			},
			orderExists: false,
			expectedErr: ErrOrderUploadedByAnotherUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(orderStorageMock)

			result, err := orderService.PutOrder(tt.args.ctx, tt.args.orderNumber)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			} else if err != nil && errors.Is(err, tt.expectedErr) {
				return
			}

			assert.Equal(t, tt.orderExists, result)
		})
	}
}

func TestOrderService_GetOrders(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	orderStorageMock := mock_storage.NewMockIOrderStorage(mockController)
	orderService := NewOrderService(orderStorageMock)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		prepareMocks   func(orderStorageMock *mock_storage.MockIOrderStorage)
		expectedResult []dto.OrderDomain
		expectedErr    error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), middleware.UserLoginContextKey, "test_user"),
			},
			prepareMocks: func(orderStorageMock *mock_storage.MockIOrderStorage) {
				orderStorageMock.
					EXPECT().
					GetOrdersByUserLogin(gomock.Any(), gomock.Any()).
					Return(
						[]dto.OrderEntity{
							{
								Id:        1,
								UserLogin: "test_user",
								Number:    "0000111122223333",
								Status:    dto.OrderStatusNew,
								Accrual: sql.NullFloat64{
									Valid: false,
								},
								UploadedAt: time.Date(2024, time.September, 10, 15, 0, 0, 0, time.UTC),
							},
							{
								Id:        2,
								UserLogin: "test_user",
								Number:    "0000111122223333",
								Status:    dto.OrderStatusProcessed,
								Accrual: sql.NullFloat64{
									Float64: testAccural,
									Valid:   true,
								},
								UploadedAt: time.Date(2024, time.September, 10, 15, 0, 0, 0, time.UTC),
							},
						},
						nil,
					)
			},
			expectedResult: []dto.OrderDomain{
				{
					Number:     "0000111122223333",
					Status:     dto.OrderStatusNew,
					Accrual:    nil,
					UploadedAt: time.Date(2024, time.September, 10, 15, 0, 0, 0, time.UTC),
				},
				{
					Number:     "0000111122223333",
					Status:     dto.OrderStatusProcessed,
					Accrual:    &testAccural,
					UploadedAt: time.Date(2024, time.September, 10, 15, 0, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks(orderStorageMock)

			result, err := orderService.GetOrders(tt.args.ctx)
			if err != nil && err != tt.expectedErr {
				require.NoError(t, err, "error not expected")
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
