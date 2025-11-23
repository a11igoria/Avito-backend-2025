package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostTeamAdd создание новой команды
func (s *Server) PostTeamAdd(ctx echo.Context) error {
	var req api.PostTeamAddJSONRequestBody

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "invalid request body",
			},
		})
	}

	// Валидация
	if req.TeamName == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
	}

	// Вызываем бизнес-логику
	team, err := s.TeamService.CreateTeam(ctx.Request().Context(), req.TeamName)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "TEAM_EXISTS",
				Message: "team already exists",
			},
		})
	}

	return ctx.JSON(http.StatusCreated, team)
}
