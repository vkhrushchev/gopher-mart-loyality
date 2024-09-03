package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

func TestOrderService_PutOrder(t *testing.T) {
	orderService := NewOrderService()

	type args struct {
		ctx         context.Context
		orderNumber string
	}

	tests := []struct {
		name        string
		args        args
		orderExists bool
		expectedErr error
	}{
		{
			name: "sucess",
			args: args{
				ctx:         context.Background(),
				orderNumber: "0123456789",
			},
			orderExists: false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := orderService.PutOrder(tt.args.ctx, tt.args.orderNumber)
			if err != nil && err != tt.expectedErr {
				require.NoError(t, err, "error not expected")
			}

			assert.Equal(t, tt.orderExists, result)
		})
	}
}

func TestOrderService_GetOrders(t *testing.T) {
	orderService := NewOrderService()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		expectedResult []dto.OrderDomain
		expectedErr    error
	}{
		{
			name: "sucess",
			args: args{
				ctx: context.Background(),
			},
			expectedResult: make([]dto.OrderDomain, 0),
			expectedErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := orderService.GetOrders(tt.args.ctx)
			if err != nil && err != tt.expectedErr {
				require.NoError(t, err, "error not expected")
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
