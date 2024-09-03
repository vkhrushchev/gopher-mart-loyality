package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

func TestAccrualService_GetAccrualInfo(t *testing.T) {
	accrualService := NewAccrualService()

	type args struct {
		ctx         context.Context
		orderNumber string
	}

	tests := []struct {
		name           string
		args           args
		expectedResult dto.AccuralInfoDomain
		expectedErr    error
	}{
		{
			name: "sucess",
			args: args{
				ctx:         context.Background(),
				orderNumber: "0123456789",
			},
			expectedResult: dto.AccuralInfoDomain{},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := accrualService.GetAccrualInfo(tt.args.ctx, tt.args.orderNumber)
			if err != nil && !errors.Is(err, tt.expectedErr) {
				require.NoError(t, err, "error not expected")
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
