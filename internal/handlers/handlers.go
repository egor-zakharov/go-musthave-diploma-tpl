package handlers

import (
	"encoding/json"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/orders"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/users"
	orders2 "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/orders"
	usersStore "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/users"
	"io"
	"net/http"
	"time"
)
import "github.com/go-chi/chi/v5"

type Server struct {
	usersSrv users.Service
	orderSrv orders.Service
}

func NewHandlers(usersSrv users.Service, orderSrv orders.Service) *Server {
	return &Server{usersSrv: usersSrv, orderSrv: orderSrv}
}

func (s *Server) Mux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware)
	r.Use(middlewares.GzipMiddleware)

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", s.register)
		r.Post("/api/user/login", s.login)
	})
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthorizedMiddleware)
		r.Post("/api/user/orders", s.createOrder)
		r.Get("/api/user/orders", s.getOrders)
	})
	return r
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}

	requestData := &dto.RegisterUserRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&requestData); err != nil {
		http.Error(w, "Incorrect input json", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Login:    requestData.Login,
		Password: requestData.Password,
	}

	if !user.IsValidLogin() {
		http.Error(w, "Invalid request login: must be presented and must be not empty", http.StatusBadRequest)
		return
	}

	if !user.IsValidPass() {
		http.Error(w, "Invalid request password: must be presented and must be not empty", http.StatusBadRequest)
		return
	}

	register, err := s.usersSrv.Register(r.Context(), user)
	if errors.Is(err, usersStore.ErrConflict) {
		http.Error(w, "User login already exists", http.StatusConflict)
		return
	}

	JWTToken, err := middlewares.BuildJWTString(register.UserID)
	if err != nil {
		http.Error(w, "Can not build auth token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: middlewares.CookieName, Value: JWTToken, Path: "/"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}

	requestData := &dto.LoginUserRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&requestData); err != nil {
		http.Error(w, "Incorrect input json", http.StatusInternalServerError)
		return
	}
	user := models.User{
		Login:    requestData.Login,
		Password: requestData.Password,
	}

	if !user.IsValidLogin() {
		http.Error(w, "Invalid request login: must be presented and must be not empty", http.StatusBadRequest)
		return
	}

	if !user.IsValidPass() {
		http.Error(w, "Invalid request password: must be presented and must be not empty", http.StatusBadRequest)
		return
	}

	login, err := s.usersSrv.Login(r.Context(), user)
	if err != nil {
		http.Error(w, "Incorrect login/password", http.StatusUnauthorized)
		return
	}

	JWTToken, err := middlewares.BuildJWTString(login.UserID)
	if err != nil {
		http.Error(w, "Can not build auth token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{Name: middlewares.CookieName, Value: JWTToken, Path: "/"})
}

func (s *Server) createOrder(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/plain" {
		http.Error(w, "Invalid request content type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderID := string(body)
	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	err = s.orderSrv.Add(r.Context(), orderID, userID)
	if errors.Is(err, orders.ErrLuhn) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if errors.Is(err, orders.ErrOrderAnotherUser) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if errors.Is(err, orders.ErrDuplicate) {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) getOrders(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value(middlewares.ContextUserIDKey).(string)
	ords, err := s.orderSrv.Get(r.Context(), userID)
	if errors.Is(err, orders2.ErrNotFound) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// заполняем модель ответа
	var resp []dto.GetOrdersResponse

	for _, order := range *ords {
		stringDate := order.UploadedAt.Format(time.RFC3339)
		date, _ := time.Parse(time.RFC3339, stringDate)

		resp = append(resp, dto.GetOrdersResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: date,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		return
	}
}
