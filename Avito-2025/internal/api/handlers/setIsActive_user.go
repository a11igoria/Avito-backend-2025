package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostUserssetIsActive изменить статус активности пользователя
func (s *Server) PostUserssetIsActive(ctx echo.Context) error {
	var req struct {
		UserID   string `json:"userid"`
		IsActive bool   `json:"isactive"`
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
	if req.UserID == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "userid is required",
			},
		})
	}

	// Проверяем существование пользователя
	user, err := s.UserService.GetUser(ctx.Request().Context(), req.UserID)
	if err != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOTFOUND",
				Message: "user not found",
			},
		})
	}

	// Обновляем статус активности
	if req.IsActive {
		s.UserService.ActivateUser(ctx.Request().Context(), req.UserID)
	} else {
		s.UserService.DeactivateUser(ctx.Request().Context(), req.UserID)
	}

	// Возвращаем обновленного пользователя
	updatedUser, _ := s.UserService.GetUser(ctx.Request().Context(), req.UserID)
	return ctx.JSON(http.StatusOK, updatedUser)
}
