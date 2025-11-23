package service

import (
	"avito-2025/internal/api"
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
func (s *TeamService) CreateTeam(ctx context.Context, name string) (*api.Team, error) {
	if name == "" {
		return nil, errors.New("team name cannot be empty")
	}

	// Создаём команду в storage
	err := s.teamRepo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	// Возвращаем в формате API
	return &api.Team{
		TeamName: name,
		Members:  []api.TeamMember{},
	}, nil
}

// GetTeamByName — получить команду по имени
func (s *TeamService) GetTeamByName(ctx context.Context, teamName string) (*api.Team, error) {
	// Получаем информацию о команде
	teamMap, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return nil, err
	}

	if teamMap == nil {
		return nil, errors.New("team not found")
	}

	// Получаем членов команды
	members, err := s.userRepo.GetActiveMembers(ctx, teamName)
	if err != nil {
		members = []map[string]interface{}{}
	}

	// Конвертируем в api.TeamMember
	apiMembers := make([]api.TeamMember, 0, len(members))
	for _, m := range members {
		apiMembers = append(apiMembers, api.TeamMember{
			UserId:   m["ID"].(string),
			Username: m["Username"].(string),
			IsActive: m["IsActive"].(bool),
		})
	}

	return &api.Team{
		TeamName: teamMap["Name"].(string),
		Members:  apiMembers,
	}, nil
}

// ListTeams — получить все команды
func (s *TeamService) ListTeams(ctx context.Context) ([]*api.Team, error) {
	teams, err := s.teamRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	if len(teams) == 0 {
		return nil, nil
	}

	result := make([]*api.Team, 0, len(teams))
	for _, t := range teams {
		teamName := t["Name"].(string)

		// Получаем членов каждой команды
		members, err := s.userRepo.GetActiveMembers(ctx, teamName)
		if err != nil {
			members = []map[string]interface{}{}
		}

		apiMembers := make([]api.TeamMember, 0, len(members))
		for _, m := range members {
			apiMembers = append(apiMembers, api.TeamMember{
				UserId:   m["ID"].(string),
				Username: m["Username"].(string),
				IsActive: m["IsActive"].(bool),
			})
		}

		result = append(result, &api.Team{
			TeamName: teamName,
			Members:  apiMembers,
		})
	}

	return result, nil
}

// UpdateTeam — обновить название команды
func (s *TeamService) UpdateTeam(ctx context.Context, oldTeamName string, newTeamName string) error {
	if newTeamName == "" {
		return errors.New("team name cannot be empty")
	}

	return s.teamRepo.Update(ctx, oldTeamName, newTeamName)
}

// DeleteTeam — удалить команду
func (s *TeamService) DeleteTeam(ctx context.Context, teamName string) error {
	return s.teamRepo.Delete(ctx, teamName)
}
