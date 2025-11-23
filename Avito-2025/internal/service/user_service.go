package service

import (
	"avito-2025/internal/api"
	"avito-2025/internal/storage"
	"context"
	"errors"
)

type UserService struct {
	userRepo *storage.UserRepository
	teamRepo *storage.TeamRepository
}

func NewUserService(userRepo *storage.UserRepository, teamRepo *storage.TeamRepository) *UserService {
	return &UserService{userRepo: userRepo, teamRepo: teamRepo}
}

// GetUser — получить пользователя по string ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*api.User, error) {
	userMap, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if userMap == nil {
		return nil, errors.New("user not found")
	}

	// Конвертируем map[string]interface{} в api.User
	return &api.User{
		UserId:   userMap["ID"].(string),
		Username: userMap["Username"].(string),
		TeamName: userMap["TeamName"].(string),
		IsActive: userMap["IsActive"].(bool),
	}, nil
}

// ActivateUser — активировать пользователя
func (s *UserService) ActivateUser(ctx context.Context, userID string) error {
	userMap, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || userMap == nil {
		return errors.New("user not found")
	}

	// Обновляем статус
	return s.userRepo.Update(ctx, userID, userMap["Username"].(string), true)
}

// DeactivateUser — деактивировать пользователя
func (s *UserService) DeactivateUser(ctx context.Context, userID string) error {
	userMap, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || userMap == nil {
		return errors.New("user not found")
	}

	// Обновляем статус
	return s.userRepo.Update(ctx, userID, userMap["Username"].(string), false)
}

// GetTeamMembers — получить активных членов команды по имени
func (s *UserService) GetTeamMembers(ctx context.Context, teamName string) ([]*api.User, error) {
	membersSlice, err := s.userRepo.GetActiveMembers(ctx, teamName)
	if err != nil {
		return nil, err
	}

	if len(membersSlice) == 0 {
		return nil, nil
	}

	// Конвертируем []map[string]interface{} в []*api.User
	result := make([]*api.User, 0, len(membersSlice))
	for _, m := range membersSlice {
		apiUser := &api.User{
			UserId:   m["ID"].(string),
			Username: m["Username"].(string),
			TeamName: m["TeamName"].(string),
			IsActive: m["IsActive"].(bool),
		}
		result = append(result, apiUser)
	}

	return result, nil
}
