package handlers

import (
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
