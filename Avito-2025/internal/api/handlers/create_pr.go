package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullRequestCreate создание pull request
func (s *Server) PostPullRequestCreate(ctx echo.Context) error {
	var req api.PostPullRequestCreateJSONBody

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
	if req.PullRequestName == "" || req.AuthorId == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "pull_request_name and author_id are required",
			},
		})
	}

	// Проверяем, что AuthorId валидный (пытаемся получить пользователя)
	author, err := s.UserService.GetUser(ctx.Request().Context(), req.AuthorId)
	if err != nil || author == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "author not found",
			},
		})
	}

	// Создаем PR
	pr, err := s.PRService.CreatePR(ctx.Request().Context(), req.PullRequestName, req.AuthorId)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "PR_EXISTS",
				Message: "PR already exists",
			},
		})
	}

	return ctx.JSON(http.StatusCreated, pr)
}
