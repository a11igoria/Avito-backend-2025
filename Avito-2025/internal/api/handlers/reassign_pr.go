package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullRequestReassign переназначение ревьювера на PR
func (s *Server) PostPullRequestReassign(ctx echo.Context) error {
	var req api.PostPullRequestReassignJSONBody

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

	// Проверяем, что PR еще не мержен
	if pr.Status == api.PullRequestStatusMERGED {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "PR_MERGED",
				Message: "cannot reassign reviewer to merged PR",
			},
		})
	}

	// Выбираем нового случайного ревьювера (старого автоматически удаляем)
	newReviewer, err := s.PRService.AssignRandomReviewer(ctx.Request().Context(), req.PullRequestId, req.OldUserId)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "NO_CANDIDATE",
				Message: "no suitable reviewer found",
			},
		})
	}

	// Возвращаем обновленный PR с новым ревьювером
	updatedPR, _ := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestId)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr":              updatedPR,
		"new_reviewer_id": newReviewer.UserId,
	})
}
