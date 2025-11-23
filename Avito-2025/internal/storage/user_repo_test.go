package storage

import (
	"avito-2025/internal/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *UserRepository {
	db, err := NewDB("postgres://postgres:postgres@localhost:5433/reviewer_service?sslmode=disable")
	require.NoError(t, err)
	return NewUserRepository(db)
}

func TestUserRepository_CRUD(t *testing.T) {
	repo := setupDB(t)
	ctx := context.Background()

	// Создание команды
	team := &domain.Team{Name: "TestTeam"}
	repoTeam := NewTeamRepository(repo.db)
	err := repoTeam.Create(ctx, team)
	require.NoError(t, err)

	// Создание
	user := &domain.User{Username: "TestUser", TeamID: 1, IsActive: true}
	err = repo.Create(ctx, user)
	require.NoError(t, err)
	require.True(t, user.ID > 0)

	// Чтение
	found, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, "TestUser", found.Username)

	// Обновление
	found.Username = "UpdatedUser"
	err = repo.Update(ctx, found)
	require.NoError(t, err)

	found2, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "UpdatedUser", found2.Username)

	// Удаление
	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	found3, err := repo.GetByID(ctx, user.ID)
	require.NoError(t, err)
	require.Nil(t, found3)
}
