package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Секретный ключ для подписи токена
var jwtKey = []byte("your_secret_key")

// Claims — структура для хранения данных в JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username, password string) (string, error) {
	// Создаем новый токен с указанием метода подписи и claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Токен будет действителен 24 часа
	})

	// Подписываем токен с использованием секретного ключа
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(username, password, inToken string) (bool, error) {
	// Парсим токен
	token, err := jwt.Parse(inToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return false, err
	}

	// Проверяем, валиден ли токен
	if !token.Valid {
		return false, fmt.Errorf("токен не валиден")
	}

	// Извлекаем claims из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("не удалось извлечь claims из токена")
	}

	// Сравниваем данные из токена с данными из запроса
	if claims["username"] != username || claims["password"] != password {
		return false, fmt.Errorf("данные в токене не совпадают с данными в запросе")
	}

	return true, nil
}
