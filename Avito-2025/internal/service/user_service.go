package service

import (
	"avito-2025/internal/domain"
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

// CreateUser — создать пользователя с проверкой команды
func (s *UserService) CreateUser(ctx context.Context, username string, teamID int) (*domain.User, error) {
	// Проверяем, существует ли команда
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	// Создаём пользователя
	user := &domain.User{
		Username: username,
		TeamID:   teamID,
		IsActive: true,
	}
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUser — получить пользователя
func (s *UserService) GetUser(ctx context.Context, userID int) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// DeactivateUser — деактивировать пользователя
func (s *UserService) DeactivateUser(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	user.IsActive = false
	return s.userRepo.Update(ctx, user)
}

// GetTeamMembers — получить активных членов команды
func (s *UserService) GetTeamMembers(ctx context.Context, teamID int) ([]*domain.User, error) {
	members, err := s.userRepo.GetActiveMembers(ctx, teamID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// DeleteUser — удалить пользователя
func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	return s.userRepo.Delete(ctx, userID)
}
