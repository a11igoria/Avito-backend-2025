package handlers

import (
	"encoding/json"
	"net/http"

	"avito-2025/internal/service"
)

type PRHandler struct {
	prService *service.PRService
}

func NewPRHandler(prService *service.PRService) *PRHandler {
	return &PRHandler{prService: prService}
}

// CreatePR — создать pull request
// POST /api/pull-requests
func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		AuthorID string `json:"author_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	pr, err := h.prService.CreatePR(r.Context(), req.Name, req.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pr)
}

// GetPR — получить PR по ID
// GET /api/pull-requests/{id}
func (h *PRHandler) GetPR(w http.ResponseWriter, r *http.Request) {
	prID := r.PathValue("id")

	pr, err := h.prService.GetPR(r.Context(), prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if pr == nil {
		http.Error(w, "PR not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pr)
}

// UpdatePRStatus — обновить статус PR
// PATCH /api/pull-requests/{id}/status
func (h *PRHandler) UpdatePRStatus(w http.ResponseWriter, r *http.Request) {
	prID := r.PathValue("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.prService.MergePR(r.Context(), prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// AssignReviewer — назначить ревьювера на PR
// POST /api/pull-requests/{id}/reviewers
func (h *PRHandler) AssignReviewer(w http.ResponseWriter, r *http.Request) {
	prID := r.PathValue("id")

	var req struct {
		ReviewerID string `json:"reviewer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.prService.AssignReviewer(r.Context(), prID, req.ReviewerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "reviewer assigned"})
}

// RemoveReviewer — удалить ревьювера с PR
// DELETE /api/pull-requests/{id}/reviewers/{reviewerId}
func (h *PRHandler) RemoveReviewer(w http.ResponseWriter, r *http.Request) {
	prID := r.PathValue("id")
	reviewerID := r.PathValue("reviewerId")

	err := h.prService.RemoveReviewer(r.Context(), prID, reviewerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "reviewer removed"})
}

// GetPRsWhereUserIsReviewer — получить все PR где пользователь ревьювер
// GET /api/users/{userId}/reviews
func (h *PRHandler) GetPRsWhereUserIsReviewer(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	prs, err := h.prService.GetPRsWhereUserIsReviewer(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prs)
}
