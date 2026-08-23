package handler

import "crypto-server/internal/service"

type Authhandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *Authhandler {
	return &Authhandler{authService: authService}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
