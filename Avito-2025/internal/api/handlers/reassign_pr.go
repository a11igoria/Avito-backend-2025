package handlers

import (
	"avito-2025/internal/api"
	"net/http"

	"github.com/labstack/echo/v4"
)

// PostPullrequestreassign переназначение ревьювера на PR
func (s *Server) PostPullrequestreassign(ctx echo.Context) error {
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

	// Проверяем, что PR еще не мержен
	if pr.Status == "MERGED" {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "PRMERGED",
				Message: "cannot reassign reviewer to merged PR",
			},
		})
	}

	// Выбираем нового случайного ревьювера (старого автоматически удаляем)
	newReviewer, err := s.PRService.AssignRandomReviewer(ctx.Request().Context(), req.PullRequestID)
	if err != nil {
		return ctx.JSON(http.StatusConflict, api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    "NOCANDIDATE",
				Message: "no suitable reviewer found",
			},
		})
	}

	// Возвращаем обновленный PR с новым ревьювером
	updatedPR, _ := s.PRService.GetPR(ctx.Request().Context(), req.PullRequestID)
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr":            updatedPR,
		"newreviewerid": newReviewer.ID,
	})
}
