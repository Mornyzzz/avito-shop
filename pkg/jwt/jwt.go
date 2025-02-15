package jwt

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

// Секретный ключ для подписи токена
var (
	jwtKey = []byte("my-avito-secret-key")
)

// Claims — структура для хранения данных в JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Срок действия
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // Время выпуска
			NotBefore: jwt.NewNumericDate(time.Now()),     // Время начала действия
			Issuer:    "your-app-name",                    // Идентификатор приложения
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok {
			ctx := context.WithValue(r.Context(), "username", claims.Username)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
