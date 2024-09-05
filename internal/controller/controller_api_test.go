package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
	mock_service "github.com/vkhrushchev/gopher-mart-loyality/internal/service/mock"
)

func TestAPIController_RegisterUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name         string
		request      *dto.APIRegisterUserRequest
		contentType  string
		setupMocks   func(userServiceMock *mock_service.MockIUserService)
		expectedCode int
	}{
		{
			name: "success",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				userServiceMock.EXPECT().
					LoginUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("test_auth_token", nil)
			},
			request: &dto.APIRegisterUserRequest{
				Login:    "test_loging",
				Password: "test_password",
			},
			contentType:  "application/json",
			expectedCode: http.StatusOK,
		},
		{
			name: "user exists",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(service.ErrUserExists)
			},
			request: &dto.APIRegisterUserRequest{
				Login:    "test_loging",
				Password: "test_password",
			},
			contentType:  "application/json",
			expectedCode: http.StatusConflict,
		},
		{
			name: "internal server error",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().
					RegisterUser(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("internal server error"))
			},
			request: &dto.APIRegisterUserRequest{
				Login:    "test_loging",
				Password: "test_password",
			},
			contentType:  "application/json",
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(userServiceMock)

			requestBytes, err := json.Marshal(tt.request)
			require.NoError(t, err, "error when marshal request")
			r := httptest.NewRequest(
				http.MethodPost,
				"/api/user/register",
				strings.NewReader(string(requestBytes)),
			)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()

			apiController.RegisterUser(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)

			if w.Result().StatusCode == http.StatusOK {
				authCookiePresent := false
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "AuthToken" {
						authCookiePresent = true
					}
				}

				assert.True(t, authCookiePresent, "no auth cookie set")
			}
		})
	}
}

func TestAPIController_LoginUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name         string
		request      *dto.APILoginUserReqest
		contentType  string
		setupMocks   func(userServiceMock *mock_service.MockIUserService)
		expectedCode int
	}{
		{
			name: "success",
			request: &dto.APILoginUserReqest{
				Login:    "test_login",
				Password: "test_password",
			},
			contentType: "application/json",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().LoginUser(gomock.Any(), gomock.Any(), gomock.Any())
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "unknown user and password",
			request: &dto.APILoginUserReqest{
				Login:    "test_login",
				Password: "test_password",
			},
			contentType: "application/json",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().LoginUser(gomock.Any(), gomock.Any(), gomock.Any()).Return("", service.ErrWrongLoginOrPassword)
			},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(userServiceMock)

			requestBytes, err := json.Marshal(tt.request)
			require.NoError(t, err, "error when marshal request")
			r := httptest.NewRequest(
				http.MethodPost,
				"/api/user/login",
				strings.NewReader(string(requestBytes)),
			)
			r.Header.Add("Content-Type", tt.contentType)

			w := httptest.NewRecorder()

			apiController.LoginUser(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)

			if w.Result().StatusCode == http.StatusOK {
				authCookiePresent := false
				cookies := w.Result().Cookies()
				for _, cookie := range cookies {
					if cookie.Name == "AuthToken" {
						authCookiePresent = true
					}
				}

				assert.True(t, authCookiePresent, "no auth cookie set")
			}
		})
	}
}

func TestAPIController_PutUserOrders(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name         string
		orderdNumber string
		setupMocks   func(orderServiceMock *mock_service.MockIOrderService)
		expectedCode int
	}{
		{
			name:         "success",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, nil)
			},
			expectedCode: http.StatusAccepted,
		},
		{
			name:         "exists",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "wrong order number",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, service.ErrOrderWrongNumber)
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "uploaded by another user",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, service.ErrOrderUploadedByAnotherUser)
			},
			expectedCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(orderServiceMock)

			r := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(tt.orderdNumber))
			w := httptest.NewRecorder()

			apiController.PutUserOrder(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}

func TestAPIController_GetUserOrders(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name        string
		c           *APIController
		epectedCode int
	}{
		{
			name:        "success",
			epectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			w := httptest.NewRecorder()

			apiController.GetUserOrders(w, r)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	}
}

func TestAPIController_GetUserBalance(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name             string
		setupMocks       func(userServiceMock *mock_service.MockIUserService)
		expectedResponse *dto.APIGetUserBalanceResponse
		expectedCode     int
	}{
		{
			name: "success",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().
					GetBalance(gomock.Any()).
					Return(
						dto.UserBalanceDomain{
							Current:  100.5,
							Withdraw: 10.5,
						},
						nil,
					)
			},
			expectedResponse: &dto.APIGetUserBalanceResponse{
				Current:   100.5,
				Withdrawn: 10.5,
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(userServiceMock)

			r := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			w := httptest.NewRecorder()

			apiController.GetUserBalance(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var apiResponse dto.APIGetUserBalanceResponse
				err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
				require.NoError(t, err, "unexpected error when parse response")
			}
		})
	}
}

func TestAPIController_WithdrawUserBalance(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name        string
		c           *APIController
		epectedCode int
	}{
		{
			name:        "success",
			epectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", nil)
			w := httptest.NewRecorder()

			apiController.WithdrawUserBalance(w, r)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	}
}

func TestAPIController_GetUserBalanaceWithdrawls(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	apiController := NewAPIController(userServiceMock, orderServiceMock)

	tests := []struct {
		name        string
		epectedCode int
	}{
		{
			name:        "success",
			epectedCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			w := httptest.NewRecorder()

			apiController.GetUserBalanaceWithdrawls(w, r)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	}
}
