package service

import (
	"avito-2025/internal/domain"
	"avito-2025/internal/storage"
	"context"
	"errors"
)

type TeamService struct {
	teamRepo *storage.TeamRepository
	userRepo *storage.UserRepository
}

func NewTeamService(teamRepo *storage.TeamRepository, userRepo *storage.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

// CreateTeam — создать команду
func (s *TeamService) CreateTeam(ctx context.Context, name string) (*domain.Team, error) {
	if name == "" {
		return nil, errors.New("team name cannot be empty")
	}

	team := &domain.Team{Name: name}
	err := s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// GetTeam — получить команду
func (s *TeamService) GetTeam(ctx context.Context, teamID int) (*domain.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}
	return team, nil
}

// UpdateTeam — обновить название команды
func (s *TeamService) UpdateTeam(ctx context.Context, teamID int, newName string) (*domain.Team, error) {
	if newName == "" {
		return nil, errors.New("team name cannot be empty")
	}

	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil || team == nil {
		return nil, errors.New("team not found")
	}

	team.Name = newName
	err = s.teamRepo.Update(ctx, team)
	if err != nil {
		return nil, err
	}
	return team, nil
}

// ListTeams — получить все команды
func (s *TeamService) ListTeams(ctx context.Context) ([]*domain.Team, error) {
	teams, err := s.teamRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return teams, nil
}

// DeleteTeam — удалить команду
func (s *TeamService) DeleteTeam(ctx context.Context, teamID int) error {
	return s.teamRepo.Delete(ctx, teamID)
}
