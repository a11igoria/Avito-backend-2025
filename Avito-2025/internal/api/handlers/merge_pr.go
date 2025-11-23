package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullrequestmerge мержирование pull request
func (s *Server) PostPullrequestmerge(ctx echo.Context) error {
	var req struct {
		PullRequestID string `json:"pullrequestid"`
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
	if req.PullRequestID == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "BAD_REQUEST",
				Message: "pullrequestid is required",
			},
		})
	}

	// Проверяем существование PR
	pr, err := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestID)
	if err != nil || pr == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOTFOUND",
				Message: "PR not found",
			},
		})
	}

	// Мержим PR
	err = s.PRService.MergePR(ctx.Request().Context(), req.PullRequestID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "MERGE_ERROR",
				Message: err.Error(),
			},
		})
	}

	// Возвращаем обновленный PR
	updatedPR, _ := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestID)
	return ctx.JSON(http.StatusOK, updatedPR)
}
