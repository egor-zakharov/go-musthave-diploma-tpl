package main

import (
	"database/sql"
	"fmt"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/clients/accrual"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/config"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/createorder"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/createwithdraw"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getbalance"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getorders"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getwithdrawals"
	loginHandle "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/login"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/registration"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/logger"
	accrualPrc "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/processors/accrual"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/server"
	balanceSrv "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/balance"
	ordersSrv "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/orders"
	usersSrv "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/users"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/balance"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/migrator"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/users"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func main() {
	cfg := config.NewConfig()
	cfg.ParseFlag()

	//Logger
	err := logger.Initialize(cfg.FlagLogLevel)
	if err != nil {
		fmt.Printf("Logger can not be initialized %s", err)
		return
	}

	//Migrator
	db, err := sql.Open("pgx", cfg.FlagDB)
	if err != nil {
		logger.Log().Sugar().Errorw("Open DB migrations crashed: ", zap.Error(err))
		panic(err)
	}
	logger.Log().Sugar().Debugw("Running DB migrations")
	newMigrator := migrator.New(db)
	err = newMigrator.Run()
	if err != nil {
		logger.Log().Sugar().Errorw("Migrations crashed with error: ", zap.Error(err))
		panic(err)
	}

	//DB
	db, err = sql.Open("pgx", cfg.FlagDB)
	if err != nil {
		logger.Log().Sugar().Errorw("Open DB storage crashed: ", zap.Error(err))
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		logger.Log().Sugar().Errorw("Cannot ping DB: ", zap.Error(err))
		panic(err)
	}
	//Storages
	usersStore := users.New(db)
	orderStore := orders.New(db)
	balanceStore := balance.New(db)
	//Services
	usersService := usersSrv.New(logger.Log(), usersStore)
	ordersService := ordersSrv.New(logger.Log(), orderStore)
	balanceService := balanceSrv.New(logger.Log(), balanceStore)
	//Clients
	accrualClient := accrual.New(logger.Log(), cfg.FlagAccAddr)
	//Processors
	accrualProc := accrualPrc.New(logger.Log(), accrualClient, orderStore, balanceStore)
	//Handlers
	registrationHandler := registration.New(usersService)
	loginHandler := loginHandle.New(usersService)
	createOrderHandler := createorder.New(ordersService)
	getOrdersHandler := getorders.New(ordersService)
	getBalanceHandler := getbalance.New(balanceService)
	createWithdrawHandler := createwithdraw.New(ordersService, balanceService)
	getWithdrawalsHandler := getwithdrawals.New(balanceService)
	//Server
	srv := server.New(registrationHandler, loginHandler, createOrderHandler, getOrdersHandler, getBalanceHandler, createWithdrawHandler, getWithdrawalsHandler)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for range ticker.C {
			accrualProc.Do()
		}
	}()

	logger.Log().Sugar().Debugw("Starting server", "address", cfg.FlagRunAddr)

	//Server
	err = http.ListenAndServe(cfg.FlagRunAddr, srv.Mux())
	if err != nil {
		logger.Log().Sugar().Errorw("Server crashed with error: ", zap.Error(err))
		panic(err)
	}

}
