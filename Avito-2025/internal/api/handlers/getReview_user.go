package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUsersgetReview получить PR, на которых пользователь ревьювер
func (s *Server) GetUsersgetReview(ctx echo.Context) error {
	userID := ctx.QueryParam("userid")

	// Валидация
	if userID == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "userid query parameter is required",
			},
		})
	}

	// Проверяем, что пользователь существует
	user, err := s.UserService.GetUser(ctx.Request().Context(), userID)
	if err != nil || user == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOTFOUND",
				Message: "user not found",
			},
		})
	}

	// Получаем PR, где этот пользователь ревьювер
	prs, err := s.PRService.GetPRsWhereUserIsReviewer(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "GET_ERROR",
				Message: err.Error(),
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"userid":       userID,
		"pullrequests": prs,
	})
}
