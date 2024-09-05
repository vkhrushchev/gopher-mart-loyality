package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
)

var orderNumberRegexp = regexp.MustCompile(`\d{16}`)

type APIController struct {
	userService  service.IUserService
	orderService service.IOrderService
}

func NewAPIController(userService service.IUserService, orderService service.IOrderService) *APIController {
	return &APIController{
		userService:  userService,
		orderService: orderService,
	}
}

func (c *APIController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	log.Infow("RegisterUser handler called.")

	if ok := checkContentType(r, w); !ok {
		return
	}

	var apiRequest dto.APIRegisterUserRequest
	if ok := parseRequest(r, w, &apiRequest); !ok {
		return
	}

	err := c.userService.RegisterUser(r.Context(), apiRequest.Login, apiRequest.Password)
	if err != nil && errors.Is(err, service.ErrUserExists) {
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
	if err != nil && errors.Is(err, service.ErrWrongLoginOrPassword) {
		log.Errorw("controller_api: unknown username or password")
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

	if ok := checkContentType(r, w); !ok {
		return
	}

	var apiRequest dto.APILoginUserReqest
	if ok := parseRequest(r, w, &apiRequest); !ok {
		return
	}

	authToken, err := c.userService.LoginUser(r.Context(), apiRequest.Login, apiRequest.Password)
	if err != nil && errors.Is(err, service.ErrWrongLoginOrPassword) {
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
	log.Infow("PutUserOrders handler called.")
	requestBodyBytes, err := io.ReadAll(r.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Errorw("controller_api: error when read request body")
	}

	orderNumber := string(requestBodyBytes)
	matched := orderNumberRegexp.Match([]byte(orderNumber))
	if !matched {
		log.Errorw(
			"controller_api: order number not matched by regexp",
		)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exist, err := c.orderService.PutOrder(r.Context(), orderNumber)
	if err != nil && errors.Is(err, service.ErrOrderWrongNumber) {
		log.Infow(
			"controller_api: wrong order number",
			"order_number", orderNumber,
		)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	} else if err != nil && errors.Is(err, service.ErrOrderUploadedByAnotherUser) {
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

	apiResponse := make([]dto.APIGetUserOrderResponseEntry, 0, len(orderDomains))
	for _, orderDomain := range orderDomains {
		apiResponseEntry := dto.APIGetUserOrderResponseEntry(orderDomain)
		apiResponse = append(apiResponse, apiResponseEntry)
	}

	json.NewEncoder(w).Encode(apiResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

	apiResponse := dto.APIGetUserBalanceResponse{
		Current:   userBalanceDomain.Current,
		Withdrawn: userBalanceDomain.Withdraw,
	}
	json.NewEncoder(w).Encode(apiResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (c *APIController) WithdrawUserBalance(w http.ResponseWriter, r *http.Request) {
	log.Infow("WithdrawUserBalance handler called.")

	w.WriteHeader(http.StatusOK)
}

func (c *APIController) GetUserBalanaceWithdrawls(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserBalanaceWithdrawls handler called.")

	w.WriteHeader(http.StatusOK)
}

func checkContentType(r *http.Request, w http.ResponseWriter) bool {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Infow(
			"controller_api: not supported \"Content-Type\" header",
			"Content-Type", contentType,
		)

		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
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
		Name:   "AuthToken",
		Value:  authToken,
		Path:   "/",
		MaxAge: 3600,
	}
	http.SetCookie(w, &authCookie)
}
