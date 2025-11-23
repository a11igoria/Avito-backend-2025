package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullrequestcreate создание pull request
func (s *Server) PostPullrequestcreate(ctx echo.Context) error {
	var req struct {
		PullRequestName string `json:"pullrequestname"`
		AuthorID        string `json:"authorid"`
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
	if req.PullRequestName == "" || req.AuthorID == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "pullrequestname and authorid are required",
			},
		})
	}

	// Проверяем, что AuthorID валидный (можно попробовать получить пользователя)
	author, err := s.UserService.GetUser(ctx.Request().Context(), req.AuthorID)
	if err != nil || author == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOTFOUND",
				Message: "author not found",
			},
		})
	}

	// Создаем PR
	pr, err := s.PRService.CreatePR(ctx.Request().Context(), req.PullRequestName, req.AuthorID)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "PREXISTS",
				Message: "PR already exists",
			},
		})
	}

	return ctx.JSON(http.StatusCreated, pr)
}
