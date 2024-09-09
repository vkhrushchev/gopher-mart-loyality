package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/controller"
	mock_service "github.com/vkhrushchev/gopher-mart-loyality/internal/service/mock"
)

var (
	runAddr      string = "localhost:8080"
	jwtSecretKey string = "jwtSecretKey"
)

func TestGopherMartLoylityApp_registerUser(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		requestBody string
		status      int
	}{
		{
			name:        "success",
			requestBody: "",
			status:      http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/user/register", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_loginUser(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		requestBody string
		status      int
	}{
		{
			name:        "success",
			requestBody: "",
			status:      http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/user/login", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_putUserOrders(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		authHeader  string
		requestBody string
		status      int
	}{
		{
			name:        "unauthrized",
			authHeader:  "",
			requestBody: "",
			status:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/user/orders", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", tt.authHeader)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_getUserOrders(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		authHeader  string
		requestBody string
		status      int
	}{
		{
			name:        "unauthrized",
			authHeader:  "",
			requestBody: "",
			status:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/api/user/orders", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", tt.authHeader)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_getUserBalance(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		authHeader  string
		requestBody string
		status      int
	}{
		{
			name:        "unauthrized",
			authHeader:  "",
			requestBody: "",
			status:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/api/user/balance", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", tt.authHeader)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_withdrawUserBalance(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		authHeader  string
		requestBody string
		status      int
	}{
		{
			name:        "unauthrized",
			authHeader:  "",
			requestBody: "",
			status:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/api/user/balance/withdraw", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", tt.authHeader)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}

func TestGopherMartLoylityApp_getUserBalanaceWithdrawls(t *testing.T) {
	t.SkipNow()

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	userServiceMock := mock_service.NewMockIUserService(mockController)
	orderServiceMock := mock_service.NewMockIOrderService(mockController)
	withdrawServiceMock := mock_service.NewMockIWithDrawService(mockController)
	apiController := controller.NewAPIController(userServiceMock, orderServiceMock, withdrawServiceMock)

	app := NewGopherMartLoylityApp(runAddr, jwtSecretKey, apiController)
	app.RegisterHandlers()

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	tests := []struct {
		name        string
		authHeader  string
		requestBody string
		status      int
	}{
		{
			name:        "unauthrized",
			authHeader:  "",
			requestBody: "",
			status:      http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, ts.URL+"/api/user/withdrawals", nil)
			require.NoError(t, err)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("Authorization", tt.authHeader)

			response, err := ts.Client().Do(request)
			require.NoError(t, err)

			assert.Equal(t, tt.status, response.StatusCode)
		})
	}
}
