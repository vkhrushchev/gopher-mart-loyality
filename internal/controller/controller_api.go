package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/dto"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
)

type APIController struct {
	userService service.IUserService
}

func NewAPIController(userService service.IUserService) *APIController {
	return &APIController{
		userService: userService,
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

func (c *APIController) PutUserOrders(w http.ResponseWriter, r *http.Request) {
	log.Infow("AddUserOrders handler called.")

	w.WriteHeader(http.StatusOK)
}

func (c *APIController) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserOrders handler called.")

	w.WriteHeader(http.StatusOK)
}

func (c *APIController) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	log.Infow("GetUserBalance handler called.")

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
