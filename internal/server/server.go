package server

import (
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/createorder"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/createwithdraw"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getbalance"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getorders"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/getwithdrawals"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/login"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/handlers/registration"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
)
import "github.com/go-chi/chi/v5"

type Server struct {
	registration   *registration.Handler
	login          *login.Handler
	createOrder    *createorder.Handler
	getOrders      *getorders.Handler
	getBalance     *getbalance.Handler
	createWithdraw *createwithdraw.Handler
	getWithdrawals *getwithdrawals.Handler
}

func New(
	registration *registration.Handler,
	login *login.Handler,
	createOrder *createorder.Handler,
	getOrders *getorders.Handler,
	getBalance *getbalance.Handler,
	createWithdraw *createwithdraw.Handler,
	getWithdrawals *getwithdrawals.Handler) *Server {
	return &Server{
		registration:   registration,
		login:          login,
		createOrder:    createOrder,
		getOrders:      getOrders,
		getBalance:     getBalance,
		createWithdraw: createWithdraw,
		getWithdrawals: getWithdrawals}
}

func (s *Server) Mux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware)
	r.Use(middlewares.GzipMiddleware)

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", s.registration.Handle)
		r.Post("/api/user/login", s.login.Handle)
	})
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthorizedMiddleware)
		r.Post("/api/user/orders", s.createOrder.Handle)
		r.Get("/api/user/orders", s.getOrders.Handle)
		r.Get("/api/user/balance", s.getBalance.Handle)
		r.Post("/api/user/balance/withdraw", s.createWithdraw.Handle)
		r.Get("/api/user/withdrawals", s.getWithdrawals.Handle)

	})
	return r
}
