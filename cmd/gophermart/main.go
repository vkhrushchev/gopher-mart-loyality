package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/app"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/controller"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/db"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

func main() {
	if err := db.ExecuteMigrations("postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable"); err != nil {
		log.Fatalw(
			"main: error when run DB migrations",
			"error", err.Error())
	}

	var sqlxdb *sqlx.DB
	var err error
	if sqlxdb, err = db.NewSqlxDB("postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable"); err != nil {
		log.Fatalw(
			"main: error when connecto to DB",
			"error", err.Error())
	}

	userStorage := storage.NewUserStorage(sqlxdb)
	orderStorage := storage.NewOrderStorage(sqlxdb)

	userService := service.NewUserService(userStorage, "salt", "jwtSecretKey")
	orderService := service.NewOrderService(orderStorage)
	withdrawService := service.NewWithdrawService()

	apiController := controller.NewAPIController(userService, orderService, withdrawService)

	app := app.NewGopherMartLoylityApp("localhost:8080", "jwtSecretKey", apiController)
	app.RegisterHandlers()

	if err := app.Run(); err != nil {
		log.Fatalw(
			"main: error when run GopherMartLoylityApp",
			"error", err.Error())
	}
}
