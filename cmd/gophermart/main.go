package main

import (
	"github.com/vkhrushchev/gopher-mart-loyality/internal/app"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/controller"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

func main() {
	userService := service.NewUserService()
	orderService := service.NewOrderService()

	apiController := controller.NewAPIController(userService, orderService)

	app := app.NewGopherMartLoylityApp("", apiController)
	app.RegisterHandlers()

	err := app.Run()
	if err != nil {
		log.Fatalw(
			"main: error when run GopherMartLoylityApp",
			"error", err.Error())
	}
}
