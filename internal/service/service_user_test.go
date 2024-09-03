package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
)

func TestUserService_RegisterUser(t *testing.T) {
	userService := NewUserService()

	type args struct {
		ctx      context.Context
		username string
		password string
	}

	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userService.RegisterUser(tt.args.ctx, tt.args.username, tt.args.password)
			if err != nil && tt.expectedErr != err {
				require.NoError(t, err, "error not expected")
			}
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	userService := NewUserService()

	type args struct {
		ctx      context.Context
		username string
		password string
	}

	tests := []struct {
		name        string
		args        args
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
				password: "test_password",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := userService.LoginUser(tt.args.ctx, tt.args.username, tt.args.password)
			if err != nil && tt.expectedErr != err {
				require.NoError(t, err, "error not expected")
			}

			assert.NotNil(t, token)
		})
	}
}

func TestUserService_GetBalance(t *testing.T) {
	userService := NewUserService()

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name           string
		args           args
		expectedResult dto.UserBalanceDomain
		expectedErr    error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			expectedResult: dto.UserBalanceDomain{},
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.GetBalance(tt.args.ctx)
			if err != nil && tt.expectedErr != err {
				require.NoError(t, err, "error not expected")
			}

			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
