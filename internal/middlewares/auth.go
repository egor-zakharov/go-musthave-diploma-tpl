package middlewares

import (
	"fmt"
	"github.com/egor-zakharov/go-musthave-diploma-tpl/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const (
	SecretKey  = "secret"
	TokenExp   = time.Hour * 3
	CookieName = "session_token"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString string) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return ""
	}

	if !token.Valid {
		logger.Log().Error("Token is not valid", zap.String("token", token.Raw))
		return ""
	}

	return claims.UserID
}

func AuthorizedMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie(CookieName)

		if err != nil || authCookie == nil || authCookie.Value == "" {
			http.Error(w, "Unauthorized requests forbidden", http.StatusUnauthorized)
			return
		}

		userID := GetUserID(authCookie.Value)

		if userID == "" {
			http.Error(w, "Unauthorized requests forbidden", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
