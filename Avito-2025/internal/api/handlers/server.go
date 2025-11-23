package handlers

import (
	"avito-2025/internal/api"
	"avito-2025/internal/service"
)

// Server реализует интерфейс ServerInterface из OpenAPI
type Server struct {
	PRService   *service.PRService
	UserService *service.UserService
	TeamService *service.TeamService
}

// NewServer конструктор
func NewServer(
	prService *service.PRService,
	userService *service.UserService,
	teamService *service.TeamService,
) *Server {
	return &Server{
		PRService:   prService,
		UserService: userService,
		TeamService: teamService,
	}
}

// ErrorResponseWithCode создает стандартный ответ об ошибке
func ErrorResponseWithCode(code string, message string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    api.ErrorResponseErrorCode(code),
			Message: message,
		},
	}
}
