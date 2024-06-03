package handlers

import (
	"encoding/json"
	"errors"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/users"
	usersStore "github.com/egor-zakharov/go-musthave-diploma-tpl/internal/storage/users"
	"net/http"
)
import "github.com/go-chi/chi/v5"

type Server struct {
	usersSrv users.Service
}

func NewHandlers(usersSrv users.Service) *Server {
	return &Server{usersSrv: usersSrv}
}

func (s *Server) Mux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middlewares.LoggerMiddleware)
	r.Use(middlewares.GzipMiddleware)

	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", s.register)
		r.Post("/api/user/login", s.login)
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
