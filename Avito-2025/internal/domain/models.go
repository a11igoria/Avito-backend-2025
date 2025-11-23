package domain

import "time"

// Team представляет команду
type Team struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// User представляет пользователя
type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	TeamID    int       `db:"team_id"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}

// PullRequest представляет PR
type PullRequest struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	AuthorID  int       `db:"author_id"`
	Status    string    `db:"status"` // например, OPEN, MERGED
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// PRReviewer — ревьювер для PR
type PRReviewer struct {
	ID         int       `db:"id"`
	PRID       int       `db:"pr_id"`
	ReviewerID int       `db:"reviewer_id"`
	AssignedAt time.Time `db:"assigned_at"`
}
