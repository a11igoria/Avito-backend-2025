package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullRequestMerge мержирование pull request
func (s *Server) PostPullRequestMerge(ctx echo.Context) error {
	var req api.PostPullRequestMergeJSONBody

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
	if req.PullRequestId == "" {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: "pull_request_id is required",
			},
		})
	}

	// Проверяем существование PR
	pr, err := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestId)
	if err != nil || pr == nil {
		return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "NOT_FOUND",
				Message: "PR not found",
			},
		})
	}

	// Мержим PR
	err = s.PRService.MergePR(ctx.Request().Context(), req.PullRequestId)
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

	// Возвращаем обновленный PR
	updatedPR, _ := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestId)
	return ctx.JSON(http.StatusOK, updatedPR)
}
