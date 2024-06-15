package login

import (
	"encoding/json"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/dto"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/middlewares"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/models"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/services/users"
	"net/http"
)

type Handler struct {
	users users.Service
}

func New(users users.Service) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
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

	login, err := h.users.Login(r.Context(), user)
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
