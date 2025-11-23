package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetTeamget получить информацию о команде по имени
func (s *Server) GetTeamget(ctx echo.Context) error {
	teamName := ctx.QueryParam("teamname")

	// Валидация
	if teamName == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "teamname query parameter is required",
			},
		})
	}

	// Получаем команду
	team, err := s.TeamService.GetTeamByName(ctx.Request().Context(), teamName)
	if err != nil || team == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOTFOUND",
				Message: "team not found",
			},
		})
	}

	return ctx.JSON(http.StatusOK, team)
}
