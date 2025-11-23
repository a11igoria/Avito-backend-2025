package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostTeamadd создание новой команды
func (s *Server) PostTeamadd(ctx echo.Context) error {
	var req struct {
		TeamName string `json:"teamname"`
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "invalid request body",
			},
		})
	}

	// Валидация
	if req.TeamName == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "teamname is required",
			},
		})
	}

	// Вызываем бизнес-логику
	team, err := s.TeamService.CreateTeam(ctx.Request().Context(), req.TeamName)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "TEAMEXISTS",
				Message: "team already exists",
			},
		})
	}

	return ctx.JSON(http.StatusCreated, team)
}
