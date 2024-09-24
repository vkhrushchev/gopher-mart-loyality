package main

import (
	"flag"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/app"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/controller"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/db"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/service"
	"github.com/vkhrushchev/gopher-mart-loyality/internal/storage"
	"go.uber.org/zap"
)

var log = zap.Must(zap.NewDevelopment()).Sugar()

type Config struct {
	runAddr              string
	databaseURI          string
	accrualSystemAddress string
}

var config Config

func main() {
	parseConfig()

	if err := db.ExecuteMigrations(config.databaseURI); err != nil {
		log.Fatalw(
			"main: error when run DB migrations",
			"error", err.Error())
	}

	var sqlxdb *sqlx.DB
	var err error
	if sqlxdb, err = db.NewSqlxDB(config.databaseURI); err != nil {
		log.Fatalw(
			"main: error when connecto to DB",
			"error", err.Error())
	}

	userStorage := storage.NewUserStorage(sqlxdb)
	orderStorage := storage.NewOrderStorage(sqlxdb)
	withdrawalStorage := storage.NewWithdrawalStorage(sqlxdb)

	userService := service.NewUserService(userStorage, "salt", "jwtSecretKey")
	orderService := service.NewOrderService(orderStorage)
	withdrawService := service.NewWithdrawalService(orderStorage, userStorage, withdrawalStorage)
	accrualService := service.NewAccrualService(config.accrualSystemAddress)
	accrualPullerService := service.NewAccrualPullerService(accrualService, orderService)

	apiController := controller.NewAPIController(userService, orderService, withdrawService, accrualPullerService)

	app := app.NewGopherMartLoylityApp(config.runAddr, "jwtSecretKey", apiController, accrualPullerService)
	app.RegisterHandlers()

	if err := app.Run(); err != nil {
		log.Fatalw(
			"main: error when run GopherMartLoylityApp",
			"error", err.Error())
	}
}

func parseConfig() {
	flag.StringVar(&config.runAddr, "a", "localhost:8080", "Run address")
	flag.StringVar(&config.databaseURI, "d", "", "Database URI")
	flag.StringVar(&config.accrualSystemAddress, "r", "", "Accural address")

	flag.Parse()

	if runAddress := os.Getenv("RUN_ADDRESS"); runAddress != "" {
		config.runAddr = runAddress
	}

	if databaseURIEnv := os.Getenv("DATABASE_URI"); databaseURIEnv != "" {
		config.databaseURI = databaseURIEnv
	}

	if accruaSystemAddressEnv := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accruaSystemAddressEnv != "" {
		config.accrualSystemAddress = accruaSystemAddressEnv
	}

	log.Infow(
		"main: config parsed",
		"runAddr", config.runAddr,
		"databaseURI", config.databaseURI,
		"accrualSystemAddress", config.accrualSystemAddress,
	)
}
