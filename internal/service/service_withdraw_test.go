package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

func TestWithdrawService_MakeWithdraw(t *testing.T) {
	withdrawService := NewWithdrawService()

	type args struct {
		ctx         context.Context
		orderNumber string
		sum         float64
	}

	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx:         context.Background(),
				orderNumber: "0123456789",
				sum:         10.5,
			},
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := withdrawService.MakeWithdraw(tt.args.ctx, tt.args.orderNumber, tt.args.sum)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}
		})
	}
}

func TestWithdrawService_GetUserWithdraws(t *testing.T) {
	withdrawService := NewWithdrawService()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		expectedResult []dto.UserWithdrawDomain
		expectedErr    error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			expectedResult: make([]dto.UserWithdrawDomain, 0),
			expectedErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := withdrawService.GetUserWithdraws(tt.args.ctx)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
