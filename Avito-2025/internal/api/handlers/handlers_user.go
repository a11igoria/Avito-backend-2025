package handlers

import (
	"encoding/json"
	"net/http"

	"avito-2025/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser — создать пользователя
// POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		TeamName string `json:"team_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Вызываем сервис (кульминация: Repo → Service → Handler)
	// TODO: Реализовать CreateUser в UserService
	// user, err := h.userService.CreateUser(r.Context(), req.Username, req.TeamName)
	// if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

// GetUser — получить пользователя по ID
// GET /api/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	// Вызываем сервис
	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetTeamMembers — получить членов команды
// GET /api/teams/{teamName}/members
func (h *UserHandler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamName := r.PathValue("teamName")

	members, err := h.userService.GetTeamMembers(r.Context(), teamName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}

// ActivateUser — активировать пользователя
// PATCH /api/users/{id}/activate
func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	err := h.userService.ActivateUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "activated"})
}

// DeactivateUser — деактивировать пользователя
// PATCH /api/users/{id}/deactivate
func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	err := h.userService.DeactivateUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deactivated"})
}
