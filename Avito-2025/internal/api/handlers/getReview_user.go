package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUsersGetReview получить PR, на которых пользователь ревьювер
func (s *Server) GetUsersGetReview(ctx echo.Context) error {
	// Получаем параметры из query
	params := api.GetUsersGetReviewParams{}
	if err := ctx.Bind(&params); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "invalid query parameters",
			},
		})
	}

	userID := params.UserId
	// Валидация
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "user_id query parameter is required",
			},
		})
	}

	// Проверяем, что пользователь существует
	user, err := s.UserService.GetUser(ctx.Request().Context(), userID)
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

	// Получаем PR, где этот пользователь ревьювер
	prs, err := s.PRService.GetPRsWhereUserIsReviewer(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prs,
	})
}
