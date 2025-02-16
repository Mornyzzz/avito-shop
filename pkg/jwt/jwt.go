package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
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

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			return
		}

		// Извлечение токена из заголовка
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			return
		}

		// Парсинг и проверка токена
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			log.Printf("Invalid token: %v", err) // Логирование ошибки
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Извлечение claims и добавление в контекст Gin
		if claims, ok := token.Claims.(*Claims); ok {
			c.Set("username", claims.Username)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Передача управления следующему обработчику
		c.Next()
	}
}
