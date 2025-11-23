package service

import (
	"avito-2025/internal/domain"
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
func (s *PRService) CreatePR(ctx context.Context, title string, authorID int) (*domain.PullRequest, error) {
	if title == "" {
		return nil, errors.New("PR title cannot be empty")
	}

	// Проверяем, существует ли автор
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil || author == nil {
		return nil, errors.New("author not found")
	}

	pr := &domain.PullRequest{
		Title:    title,
		AuthorID: authorID,
		Status:   "OPEN",
	}
	err = s.prRepo.Create(ctx, pr)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

// GetPR — получить PR
func (s *PRService) GetPR(ctx context.Context, prID int) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, errors.New("PR not found")
	}
	return pr, nil
}

// AssignReviewerToPR — назначить ревьювера на PR
func (s *PRService) AssignReviewerToPR(ctx context.Context, prID int, reviewerID int) error {
	// Проверяем, существует ли PR
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return errors.New("PR not found")
	}

	// Проверяем статус PR (нельзя назначать ревьюверов на merged PR)
	if pr.Status == "MERGED" {
		return errors.New("cannot assign reviewers to merged PR")
	}

	// Проверяем, существует ли ревьювер
	reviewer, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil || reviewer == nil {
		return errors.New("reviewer not found")
	}

	// Нельзя назначать автора на ревью своего PR
	if pr.AuthorID == reviewerID {
		return errors.New("author cannot review their own PR")
	}

	// Нельзя назначать неактивного пользователя
	if !reviewer.IsActive {
		return errors.New("reviewer is not active")
	}

	// Назначаем ревьювера
	return s.prReviewerRepo.AssignReviewer(ctx, prID, reviewerID)
}

// MergePR — мержить PR
func (s *PRService) MergePR(ctx context.Context, prID int) error {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return errors.New("PR not found")
	}

	if pr.Status == "MERGED" {
		return errors.New("PR already merged")
	}

	return s.prRepo.UpdateStatus(ctx, prID, "MERGED")
}

// GetPRReviewers — получить ревьюверов PR
func (s *PRService) GetPRReviewers(ctx context.Context, prID int) ([]*domain.PRReviewer, error) {
	reviewers, err := s.prReviewerRepo.ListByPR(ctx, prID)
	if err != nil {
		return nil, err
	}
	return reviewers, nil
}

// RemoveReviewerFromPR — удалить ревьювера с PR
func (s *PRService) RemoveReviewerFromPR(ctx context.Context, prID int, reviewerID int) error {
	return s.prReviewerRepo.RemoveReviewer(ctx, prID, reviewerID)
}

// ListPRsByAuthor — получить все PR автора
func (s *PRService) ListPRsByAuthor(ctx context.Context, authorID int) ([]*domain.PullRequest, error) {
	prs, err := s.prRepo.GetByAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// ListOpenPRs — получить все открытые PR
func (s *PRService) ListOpenPRs(ctx context.Context) ([]*domain.PullRequest, error) {
	prs, err := s.prRepo.GetByStatus(ctx, "OPEN")
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func (s *PRService) AssignRandomReviewer(ctx context.Context, prID int) (*domain.User, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil || pr == nil {
		return nil, errors.New("PR not found")
	}

	// Получаем всех активных пользователей из команды автора
	author, err := s.userRepo.GetByID(ctx, pr.AuthorID)
	if err != nil || author == nil {
		return nil, errors.New("Author not found")
	}
	members, err := s.userRepo.GetActiveMembers(ctx, author.TeamID)
	if err != nil {
		return nil, err
	}

	// Исключаем автора из списка кандидатов
	candidates := make([]*domain.User, 0)
	for _, u := range members {
		if u.ID != pr.AuthorID {
			candidates = append(candidates, u)
		}
	}
	if len(candidates) == 0 {
		return nil, errors.New("No available reviewers in team")
	}

	// Случайный выбор из кандидатов
	rand.Seed(time.Now().UnixNano())
	chosen := candidates[rand.Intn(len(candidates))]

	// Назначаем ревьювера через репозиторий
	err = s.prReviewerRepo.AssignReviewer(ctx, prID, chosen.ID)
	if err != nil {
		return nil, err
	}
	return chosen, nil
}
