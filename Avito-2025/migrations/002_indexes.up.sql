-- Индексы для производительности
CREATE INDEX idx_users_team_id ON users(team_id);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX idx_pr_status ON pull_requests(status);
CREATE INDEX idx_pr_reviewers_pr_id ON pr_reviewers(pr_id);
CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);