-- Создание таблицы команд
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы pull requests
CREATE TABLE pull_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(512) NOT NULL,
    author_id INT NOT NULL REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы ревьюверов PR
CREATE TABLE pr_reviewers (
    id SERIAL PRIMARY KEY,
    pr_id INT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id INT NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(pr_id, reviewer_id)
);
