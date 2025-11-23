package service

import (
	"context"
	"testing"

	"avito-2025/internal/domain"

	"github.com/stretchr/testify/require"
)

// Мок UserRepo для теста
type MockUserRepo struct {
	users     map[int]*domain.User
	teamUsers map[int][]*domain.User
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return m.users[id], nil
}
func (m *MockUserRepo) GetActiveMembers(ctx context.Context, teamID int) ([]*domain.User, error) {
	return m.teamUsers[teamID], nil
}

// Мок PRRepo для теста
type MockPRRepo struct {
	prs map[int]*domain.PullRequest
}

func (m *MockPRRepo) GetByID(ctx context.Context, id int) (*domain.PullRequest, error) {
	return m.prs[id], nil
}

// Мок PRReviewerRepo для теста
type MockPRReviewerRepo struct {
	assigned map[int]int
}

func (m *MockPRReviewerRepo) AssignReviewer(ctx context.Context, prID int, reviewerID int) error {
	if m.assigned == nil {
		m.assigned = make(map[int]int)
	}
	m.assigned[prID] = reviewerID
	return nil
}

type PRRepoInterface interface {
	GetByID(ctx context.Context, id int) (*domain.PullRequest, error)
}
type PRReviewerRepoInterface interface {
	AssignReviewer(ctx context.Context, prID int, reviewerID int) error
}
type UserRepoInterface interface {
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetActiveMembers(ctx context.Context, teamID int) ([]*domain.User, error)
}

type PRService struct {
	prRepo         PRRepoInterface
	prReviewerRepo PRReviewerRepoInterface
	userRepo       UserRepoInterface
}

// Юнит-тест AssignRandomReviewer
func TestAssignRandomReviewer(t *testing.T) {
	// Моковые данные
	mockUsers := map[int]*domain.User{
		1: {ID: 1, Username: "author", TeamID: 100, IsActive: true},
		2: {ID: 2, Username: "reviewer", TeamID: 100, IsActive: true},
	}
	mockTeamMembers := map[int][]*domain.User{
		100: {
			{ID: 1, Username: "author", TeamID: 100, IsActive: true},
			{ID: 2, Username: "reviewer", TeamID: 100, IsActive: true},
		},
	}
	mockPRs := map[int]*domain.PullRequest{
		10: {ID: 10, Title: "test PR", AuthorID: 1, Status: "OPEN"},
	}

	userRepo := &MockUserRepo{users: mockUsers, teamUsers: mockTeamMembers}
	prRepo := &MockPRRepo{prs: mockPRs}
	prReviewerRepo := &MockPRReviewerRepo{}

	prService := &PRService{
		prRepo:         prRepo,
		prReviewerRepo: prReviewerRepo,
		userRepo:       userRepo,
	}

	// Вызов бизнес-метода
	reviewer, err := prService.AssignRandomReviewer(context.Background(), 10)

	require.NoError(t, err)
	require.NotNil(t, reviewer)
	require.NotEqual(t, reviewer.ID, 1)                        // Автор не может быть выбран
	require.Equal(t, prReviewerRepo.assigned[10], reviewer.ID) // Проверяем назначение
}
