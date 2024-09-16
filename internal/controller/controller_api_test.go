package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
	mock_service "github.com/vkhrushchev/gopher-mart-loyality/internal/service/mock"
)

var (
	testAccural float64 = 10.5
)

func TestAPIController_RegisterUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name         string
		request      dto.APIRegisterUserRequest
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
			request: dto.APIRegisterUserRequest{
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
			request: dto.APIRegisterUserRequest{
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
			request: dto.APIRegisterUserRequest{
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
			defer w.Result().Body.Close()

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
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name         string
		request      dto.APILoginUserRequest
		contentType  string
		setupMocks   func(userServiceMock *mock_service.MockIUserService)
		expectedCode int
	}{
		{
			name: "success",
			request: dto.APILoginUserRequest{
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
			request: dto.APILoginUserRequest{
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
			defer w.Result().Body.Close()

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

func TestAPIController_PutOrder(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name         string
		orderdNumber string
		setupMocks   func(orderServiceMock *mock_service.MockIOrderService, accrualPullerServiceMock *mock_service.MockIAccrualPullerService)
		expectedCode int
	}{
		{
			name:         "success",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService, accrualPullerServiceMock *mock_service.MockIAccrualPullerService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, nil)

				accrualPullerServiceMock.EXPECT().AddGetAccrualInfoTask(gomock.Any(), gomock.Any())
			},
			expectedCode: http.StatusAccepted,
		},
		{
			name:         "exists",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService, accrualPullerServiceMock *mock_service.MockIAccrualPullerService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(true, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "wrong order number",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService, accrualPullerServiceMock *mock_service.MockIAccrualPullerService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, service.ErrOrderWrongNumber)
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name:         "uploaded by another user",
			orderdNumber: "1111222233334444",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService, accrualPullerServiceMock *mock_service.MockIAccrualPullerService) {
				orderServiceMock.EXPECT().
					PutOrder(gomock.Any(), gomock.Any()).
					Return(false, service.ErrOrderUploadedByAnotherUser)
			},
			expectedCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(orderServiceMock, accrualPullerServiceMock)

			r := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader(tt.orderdNumber))
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()

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
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name                string
		setupMocks          func(orderServiceMock *mock_service.MockIOrderService)
		expectedCode        int
		expectedAPIResponse []dto.APIOrderResponse
	}{
		{
			name: "success",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					GetOrders(gomock.Any()).
					Return(
						[]dto.OrderDomain{
							{
								Number:     "1111222233334444",
								Status:     dto.OrderStatusNew,
								Accrual:    &testAccural,
								UploadedAt: time.Date(2024, time.September, 5, 10, 0, 0, 0, time.UTC),
							},
						},
						nil,
					)
			},
			expectedCode: http.StatusOK,
			expectedAPIResponse: []dto.APIOrderResponse{
				{
					Number:     "1111222233334444",
					Status:     dto.OrderStatusNew,
					Accrual:    &testAccural,
					UploadedAt: time.Date(2024, time.September, 5, 10, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "orders not found",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					GetOrders(gomock.Any()).
					Return(
						[]dto.OrderDomain{},
						nil,
					)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "internal server error",
			setupMocks: func(orderServiceMock *mock_service.MockIOrderService) {
				orderServiceMock.EXPECT().
					GetOrders(gomock.Any()).
					Return(
						nil,
						errors.New("internal server error"),
					)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(orderServiceMock)

			r := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()

			apiController.GetUserOrders(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var apiResponse []dto.APIOrderResponse
				err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
				require.NoError(t, err, "unexpected error when parse response")

				assert.Equal(t, tt.expectedAPIResponse, apiResponse)
			}
		})
	}
}

func TestAPIController_GetUserBalance(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name             string
		setupMocks       func(userServiceMock *mock_service.MockIUserService)
		expectedResponse *dto.APIUserBalance
		expectedCode     int
	}{
		{
			name: "success",
			setupMocks: func(userServiceMock *mock_service.MockIUserService) {
				userServiceMock.EXPECT().
					GetBalance(gomock.Any()).
					Return(
						dto.UserBalanceDomain{
							Current:    100.5,
							Withdrawal: 10.5,
						},
						nil,
					)
			},
			expectedResponse: &dto.APIUserBalance{
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
			defer w.Result().Body.Close()

			apiController.GetUserBalance(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var apiResponse dto.APIUserBalance
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
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name         string
		apiRequest   dto.APIPutOrderWithdrawnRequest
		setupMocks   func(withdrawService *mock_service.MockIWithdrawalService)
		expectedCode int
	}{
		{
			name: "success",
			apiRequest: dto.APIPutOrderWithdrawnRequest{
				Order: "1111222233334444",
				Sum:   10.5,
			},
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					DoWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "no funds on balance",
			apiRequest: dto.APIPutOrderWithdrawnRequest{
				Order: "1111222233334444",
				Sum:   10.5,
			},
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					DoWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(service.ErrNoFundsOnBalance)
			},
			expectedCode: http.StatusPaymentRequired,
		},
		{
			name: "wrong order number",
			apiRequest: dto.APIPutOrderWithdrawnRequest{
				Order: "1111222233334444",
				Sum:   10.5,
			},
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					DoWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(service.ErrOrderWrongNumber)
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "internal server error",
			apiRequest: dto.APIPutOrderWithdrawnRequest{
				Order: "1111222233334444",
				Sum:   10.5,
			},
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					DoWithdrawal(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("internal server error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(withdrawalServiceMock)

			requestBytes, err := json.Marshal(tt.apiRequest)
			require.NoError(t, err, "error when marshal request")
			r := httptest.NewRequest(
				http.MethodPost,
				"/api/user/balance/withdraw",
				strings.NewReader(string(requestBytes)),
			)
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()

			apiController.WithdrawUserBalance(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}

func TestAPIController_GetUserBalanaceWithdrawls(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawalServiceMock := mock_service.NewMockIWithdrawalService(mockController)
	accrualPullerServiceMock := mock_service.NewMockIAccrualPullerService(mockController)
	apiController := NewAPIController(
		userServiceMock,
		orderServiceMock,
		withdrawalServiceMock,
		accrualPullerServiceMock,
	)

	tests := []struct {
		name                string
		setupMocks          func(withdrawService *mock_service.MockIWithdrawalService)
		expectedAPIResponse []dto.APIOrderWithdrawn
		expectedCode        int
	}{
		{
			name: "success",
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					GetUserWithdrawals(gomock.Any()).
					Return(
						[]dto.OrderWithdrawalDomain{
							{
								OrderNumber:   "1111222233334444",
								WithdrawalSum: 10.5,
								ProcessedAt:   time.Date(2024, time.September, 5, 10, 0, 0, 0, time.UTC),
							},
						},
						nil,
					)
			},
			expectedAPIResponse: []dto.APIOrderWithdrawn{
				{
					OrderNumber:   "1111222233334444",
					WithdrawalSum: 10.5,
					ProcessedAt:   time.Date(2024, time.September, 5, 10, 0, 0, 0, time.UTC),
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "withdraws not found",
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					GetUserWithdrawals(gomock.Any()).
					Return(
						[]dto.OrderWithdrawalDomain{},
						nil,
					)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "internal server error",
			setupMocks: func(withdrawService *mock_service.MockIWithdrawalService) {
				withdrawService.EXPECT().
					GetUserWithdrawals(gomock.Any()).
					Return(
						nil,
						errors.New("internal server error"),
					)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(withdrawalServiceMock)

			r := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			w := httptest.NewRecorder()
			defer w.Result().Body.Close()

			apiController.GetUserBalanaceWithdrawals(w, r)

			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
			if w.Result().StatusCode == http.StatusOK {
				var apiResponse []dto.APIOrderWithdrawn
				err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
				require.NoError(t, err, "unexpected error when parse response")
			}
		})
	}
}
