package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetTeamGet получить информацию о команде по имени
func (s *Server) GetTeamGet(ctx echo.Context, params api.GetTeamGetParams) error {
	teamName := params.TeamName

	// Валидация
	if teamName == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "team_name query parameter is required",
			},
		})
	}

	// Получаем команду
	team, err := s.TeamService.GetTeamByName(ctx.Request().Context(), teamName)
	if err != nil || team == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "team not found",
			},
		})
	}

	return ctx.JSON(http.StatusOK, team)
}
