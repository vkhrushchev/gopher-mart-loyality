package app

import (
	"net/http"

	"github.com/vkhrushchev/gopher-mart-loyality/internal/controller"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/middleware"

	"github.com/go-chi/chi/v5"

	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

type GopherMartLoylityApp struct {
	apiController *controller.APIController
	router        chi.Router
	runAddr       string
	jwtSecretKey  string
}

func NewGopherMartLoylityApp(
	runAddr string,
	jwtSecretKey string,
	apiController *controller.APIController) *GopherMartLoylityApp {
	return &GopherMartLoylityApp{
		apiController: apiController,
		router:        chi.NewRouter(),
		runAddr:       runAddr,
		jwtSecretKey:  jwtSecretKey,
	}
}

func (a *GopherMartLoylityApp) RegisterHandlers() {
	a.router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", a.apiController.RegisterUser)
		r.Post("/login", a.apiController.LoginUser)

		r.Group(func(r chi.Router) {
			r.Use(middleware.NewJWTAuthMiddleware(a.jwtSecretKey))

			r.Post("/orders", a.apiController.PutUserOrder)
			r.Get("/orders", a.apiController.GetUserOrders)
			r.Get("/balance", a.apiController.GetUserBalance)
			r.Post("/balance/withdraw", a.apiController.WithdrawUserBalance)
			r.Get("/withdrawals", a.apiController.GetUserBalanaceWithdrawals)
		})
	})
}

func (a *GopherMartLoylityApp) Run() error {
	log.Infow(
		"app: GopherMartLoylityApp stating",
		"runAddr", a.runAddr,
	)

	err := http.ListenAndServe(a.runAddr, a.router)
	if err != nil {
		return err
	}

	return nil
}
