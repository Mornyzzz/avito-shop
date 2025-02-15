package entity

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string
}

// AuthRequest представляет запрос на аутентификацию
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse представляет ответ на успешную аутентификацию
type AuthResponse struct {
	Token string `json:"token"` // JWT-токен
}
