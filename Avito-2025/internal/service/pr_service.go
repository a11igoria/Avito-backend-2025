package service

import (
	"avito-2025/internal/api"
	"avito-2025/internal/storage"
	"context"
	"errors"
	"math/rand"
	"time"
)

type PRService struct {
	prRepo         *storage.PRRepository
	prReviewerRepo *storage.PRReviewerRepository
	userRepo       *storage.UserRepository
}

func NewPRService(prRepo *storage.PRRepository, prReviewerRepo *storage.PRReviewerRepository, userRepo *storage.UserRepository) *PRService {
	return &PRService{prRepo: prRepo, prReviewerRepo: prReviewerRepo, userRepo: userRepo}
}

// CreatePR — создать pull request
func (s *PRService) CreatePR(ctx context.Context, name string, authorID string) (*api.PullRequest, error) {
	if name == "" {
		return nil, errors.New("PR name cannot be empty")
	}

	// Проверяем, существует ли автор
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, errors.New("author not found")
	}

	// Создаём PR
	prID, err := s.prRepo.Create(ctx, name, authorID, string(api.PullRequestStatusOPEN))
	if err != nil {
		return nil, err
	}

	// Возвращаем в формате API
	return &api.PullRequest{
		PullRequestId:     prID,
		PullRequestName:   name,
		AuthorId:          authorID,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: []string{},
	}, nil
}

// GetPR — получить PR
func (s *PRService) GetPR(ctx context.Context, prID string) (*api.PullRequest, error) {
	prMap, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}

	if prMap == nil {
		return nil, errors.New("PR not found")
	}

	// Получаем ревьюверов
	reviewers, err := s.prReviewerRepo.GetByPR(ctx, prID)
	if err != nil {
		reviewers = []string{}
	}

	// Преобразуем статус
	status := api.PullRequestStatusOPEN
	if prMap["Status"].(string) == string(api.PullRequestStatusMERGED) {
		status = api.PullRequestStatusMERGED
	}

	return &api.PullRequest{
		PullRequestId:     prMap["ID"].(string),
		PullRequestName:   prMap["Name"].(string),
		AuthorId:          prMap["AuthorID"].(string),
		Status:            status,
		AssignedReviewers: reviewers,
	}, nil
}

// MergePR — мержить PR
func (s *PRService) MergePR(ctx context.Context, prID string) error {
	prMap, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || prMap == nil {
		return errors.New("PR not found")
	}

	// Проверяем, уже ли PR мержен
	if prMap["Status"].(string) == string(api.PullRequestStatusMERGED) {
		return errors.New("PR already merged")
	}

	// Обновляем статус на MERGED
	return s.prRepo.UpdateStatus(ctx, prID, string(api.PullRequestStatusMERGED))
}

// GetPRsWhereUserIsReviewer — получить все PR где юзер ревьювер
func (s *PRService) GetPRsWhereUserIsReviewer(ctx context.Context, userID string) ([]api.PullRequest, error) {
	prIDs, err := s.prReviewerRepo.GetPRsByReviewer(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]api.PullRequest, 0, len(prIDs))
	for _, prID := range prIDs {
		pr, err := s.GetPR(ctx, prID)
		if err != nil {
			continue
		}
		result = append(result, *pr)
	}

	return result, nil
}

// AssignRandomReviewer — назначить случайного ревьювера
func (s *PRService) AssignRandomReviewer(ctx context.Context, prID string, oldUserID string) (*api.TeamMember, error) {
	prMap, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || prMap == nil {
		return nil, errors.New("PR not found")
	}

	// Проверяем что PR еще не мержен
	if prMap["Status"].(string) == string(api.PullRequestStatusMERGED) {
		return nil, errors.New("cannot reassign reviewer on merged PR")
	}

	// Получаем автора
	authorID := prMap["AuthorID"].(string)
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, errors.New("author not found")
	}

	// Получаем активных членов команды
	members, err := s.userRepo.GetActiveMembers(ctx, author["TeamName"].(string))
	if err != nil {
		return nil, err
	}

	// Фильтруем: исключаем автора и старого ревьювера
	candidates := make([]map[string]interface{}, 0)
	for _, m := range members {
		mID := m["ID"].(string)
		mIsActive := m["IsActive"].(bool)
		if mID != authorID && mID != oldUserID && mIsActive {
			candidates = append(candidates, m)
		}
	}

	if len(candidates) == 0 {
		return nil, errors.New("no available reviewers")
	}

	// Случайный выбор
	rand.Seed(time.Now().UnixNano())
	chosen := candidates[rand.Intn(len(candidates))]

	// Удаляем старого ревьювера
	if oldUserID != "" {
		_ = s.prReviewerRepo.RemoveReviewer(ctx, prID, oldUserID)
	}

	// Назначаем нового
	chosenID := chosen["ID"].(string)
	err = s.prReviewerRepo.AssignReviewer(ctx, prID, chosenID)
	if err != nil {
		return nil, err
	}

	return &api.TeamMember{
		UserId:   chosenID,
		Username: chosen["Username"].(string),
		IsActive: chosen["IsActive"].(bool),
	}, nil
}

// AssignReviewer — назначить ревьювера на PR
func (s *PRService) AssignReviewer(ctx context.Context, prID string, reviewerID string) error {
	// Проверяем, существует ли PR
	prMap, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || prMap == nil {
		return errors.New("PR not found")
	}

	// Проверяем статус PR (нельзя назначать ревьюверов на merged PR)
	if prMap["Status"].(string) == string(api.PullRequestStatusMERGED) {
		return errors.New("cannot assign reviewers to merged PR")
	}

	// Проверяем, существует ли ревьювер
	reviewer, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil || reviewer == nil {
		return errors.New("reviewer not found")
	}

	// Нельзя назначать автора на ревью своего PR
	if prMap["AuthorID"].(string) == reviewerID {
		return errors.New("author cannot review their own PR")
	}

	// Нельзя назначать неактивного пользователя
	if !reviewer["IsActive"].(bool) {
		return errors.New("reviewer is not active")
	}

	// Назначаем ревьювера
	return s.prReviewerRepo.AssignReviewer(ctx, prID, reviewerID)
}

// RemoveReviewer — удалить ревьювера с PR
func (s *PRService) RemoveReviewer(ctx context.Context, prID string, reviewerID string) error {
	return s.prReviewerRepo.RemoveReviewer(ctx, prID, reviewerID)
}
