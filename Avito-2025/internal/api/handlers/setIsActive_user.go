package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostUsersSetIsActive изменить статус активности пользователя
func (s *Server) PostUsersSetIsActive(ctx echo.Context) error {
	var req api.PostUsersSetIsActiveJSONBody

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
	if req.UserId == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "user_id is required",
			},
		})
	}

	// Проверяем существование пользователя
	user, err := s.UserService.GetUser(ctx.Request().Context(), req.UserId)
	if err != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "user not found",
			},
		})
	}

	// Обновляем статус активности
	if req.IsActive {
		s.UserService.ActivateUser(ctx.Request().Context(), req.UserId)
	} else {
		s.UserService.DeactivateUser(ctx.Request().Context(), req.UserId)
	}

	// Возвращаем обновленного пользователя
	updatedUser, _ := s.UserService.GetUser(ctx.Request().Context(), req.UserId)
	return ctx.JSON(http.StatusOK, updatedUser)
}
