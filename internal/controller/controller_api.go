package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
)

type APIController struct {
	userService          service.IUserService
	orderService         service.IOrderService
	withdrawService      service.IWithdrawalService
	accrualPullerService service.IAccrualPullerService
}

func NewAPIController(
	userService service.IUserService,
	orderService service.IOrderService,
	withdrawService service.IWithdrawalService,
	accrualPullerService service.IAccrualPullerService) *APIController {
	return &APIController{
		userService:          userService,
		orderService:         orderService,
		withdrawService:      withdrawService,
		accrualPullerService: accrualPullerService,
	}
}

func (c *APIController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Infow("RegisterUser handler called.")

	var apiRequest dto.APIRegisterUserRequest
	if ok := parseRequest(r, w, &apiRequest); !ok {
		return
	}

	err := c.userService.RegisterUser(r.Context(), apiRequest.Login, apiRequest.Password)
	if errors.Is(err, service.ErrUserExists) {
		log.Infow(
			"controller_api: user exists",
			"username", apiRequest.Login,
		)

		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Errorw(
			"controller_api: unexpected error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authToken, err := c.userService.LoginUser(r.Context(), apiRequest.Login, apiRequest.Password)
	if errors.Is(err, service.ErrWrongLoginOrPassword) {
		log.Errorw("controller_api: unknown username or password")

		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Errorw(
			"controller_api: unexpected error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setAuthTokenCookie(w, authToken)
	w.WriteHeader(http.StatusOK)
}

func (c *APIController) LoginUser(w http.ResponseWriter, r *http.Request) {
	log.Infow("LoginUser handler called.")

	var apiRequest dto.APILoginUserRequest
	if ok := parseRequest(r, w, &apiRequest); !ok {
		return
	}

	authToken, err := c.userService.LoginUser(r.Context(), apiRequest.Login, apiRequest.Password)
	if errors.Is(err, service.ErrWrongLoginOrPassword) {
		log.Errorw("controller_api: unknown username or password")

		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Errorw(
			"controller_api: unexpected error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setAuthTokenCookie(w, authToken)
	w.WriteHeader(http.StatusOK)
}

func (c *APIController) PutUserOrder(w http.ResponseWriter, r *http.Request) {
	log.Infow("controller_api: PutUserOrders handler called")
	requestBodyBytes, err := io.ReadAll(r.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Errorw("controller_api: error when read request body")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	orderNumber := string(requestBodyBytes)
	log.Infow(
		"controller_api: start processing order",
		"order_number", orderNumber,
	)

	if _, err := strconv.Atoi(orderNumber); err != nil {
		log.Infow(
			"controller_api: order number not matched by regexp",
			"order_number", orderNumber,
		)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exist, err := c.orderService.PutOrder(r.Context(), orderNumber)
	if errors.Is(err, service.ErrOrderWrongNumber) {
		log.Infow(
			"controller_api: wrong order number",
			"order_number", orderNumber,
		)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	} else if errors.Is(err, service.ErrOrderUploadedByAnotherUser) {
		log.Infow(
			"controller_api: order uploaded by another user",
			"order_number", orderNumber,
		)

		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Infow(
			"controller_api: unexpected internal servcer error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exist {
		log.Infow(
			"controller_api: order already uploaded",
			"order_number", orderNumber,
		)

		w.WriteHeader(http.StatusOK)
		return
	}

	c.accrualPullerService.AddGetAccrualInfoTask(orderNumber)

	w.WriteHeader(http.StatusAccepted)
}

func (c *APIController) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserOrders handler called.")

	orderDomains, err := c.orderService.GetOrders(r.Context())
	if err != nil {
		log.Errorw(
			"controller_api: unexpected internal server error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if orderDomains != nil && len(orderDomains) == 0 {
		log.Infow("controller_api: no orders found for user")

		w.WriteHeader(http.StatusNoContent)
		return
	}

	apiResponse := make([]dto.APIOrderResponse, 0, len(orderDomains))
	for _, orderDomain := range orderDomains {
		apiResponseEntry := dto.APIOrderResponse(orderDomain)
		apiResponse = append(apiResponse, apiResponseEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)
}

func (c *APIController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserBalance handler called.")

	userBalanceDomain, err := c.userService.GetBalance(r.Context())
	if err != nil {
		log.Errorw(
			"controller_api: get user balance failed",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiResponse := dto.APIUserBalance{
		Current:   userBalanceDomain.Current,
		Withdrawn: userBalanceDomain.Withdrawal,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)
}

func (c *APIController) WithdrawUserBalance(w http.ResponseWriter, r *http.Request) {
	log.Infow("WithdrawUserBalance handler called.")

	var apiRequest dto.APIPutOrderWithdrawnRequest
	if ok := parseRequest(r, w, &apiRequest); !ok {
		return
	}

	err := c.withdrawService.DoWithdrawal(r.Context(), apiRequest.Order, apiRequest.Sum)
	if errors.Is(err, service.ErrNoFundsOnBalance) {
		log.Infow("controller_api: no funds on balance")

		w.WriteHeader(http.StatusPaymentRequired)
		return
	} else if errors.Is(err, service.ErrOrderWrongNumber) {
		log.Infow(
			"controller_api: wrong order number",
			"order_number", apiRequest.Order,
		)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	} else if err != nil {
		log.Errorw(
			"controller_api: internal server error",
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *APIController) GetUserBalanaceWithdrawals(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserBalanaceWithdrawals handler called.")

	userWithdrawDomains, err := c.withdrawService.GetUserWithdrawals(r.Context())
	if err != nil {
		log.Infow(
			"controller_api: internal server error",
			"error", err.Error(),
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(userWithdrawDomains) == 0 {
		log.Infow("controller_api: no withdraws found")

		w.WriteHeader(http.StatusNoContent)
		return
	}

	apiResponse := make([]dto.APIOrderWithdrawn, 0, len(userWithdrawDomains))
	for _, userWithdraw := range userWithdrawDomains {
		apiResponseEntry := dto.APIOrderWithdrawn(userWithdraw)
		apiResponse = append(apiResponse, apiResponseEntry)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiResponse)
}

func parseRequest(r *http.Request, w http.ResponseWriter, apiRequest any) bool {
	if err := json.NewDecoder(r.Body).Decode(apiRequest); err != nil {
		log.Errorw(
			"controller_api: error when decode request body from json",
			"erorr", err.Error(),
		)

		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func setAuthTokenCookie(w http.ResponseWriter, authToken string) {
	authCookie := http.Cookie{
		Name:   middleware.AuthTokenCoockieName,
		Value:  authToken,
		Path:   "/",
		MaxAge: 3600,
	}
	http.SetCookie(w, &authCookie)
}
