-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    external_id TEXT UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE users IS 'Все пользователи системы';
COMMENT ON COLUMN users.id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN users.username IS 'Уникальный логин пользователя в системе';
COMMENT ON COLUMN users.email IS 'Уникальный email пользователя';
COMMENT ON COLUMN users.is_active IS 'Флаг активности пользователя (false - не может быть назначен на ревью)';
COMMENT ON COLUMN users.created_at IS 'Время создания пользователя';
COMMENT ON COLUMN users.updated_at IS 'Время последнего обновления пользователя';

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY,
    external_id TEXT UNIQUE NOT NULL,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

COMMENT ON TABLE teams IS 'Команды разработчиков';
COMMENT ON COLUMN teams.id IS 'Уникальный идентификатор команды';
COMMENT ON COLUMN teams.name IS 'Уникальное название команды';
COMMENT ON COLUMN teams.description IS 'Описание команды (опционально)';
COMMENT ON COLUMN teams.created_at IS 'Время создания команды';
COMMENT ON COLUMN teams.updated_at IS 'Время последнего обновления команды';

CREATE TABLE IF NOT EXISTS team_members (
    id UUID PRIMARY KEY,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'MEMBER',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(team_id, user_id)
);

COMMENT ON TABLE team_members IS 'Состав команд (связь пользователей с командами)';
COMMENT ON COLUMN team_members.id IS 'Уникальный идентификатор членства в команде';
COMMENT ON COLUMN team_members.team_id IS 'Идентификатор команды';
COMMENT ON COLUMN team_members.user_id IS 'Идентификатор пользователя';
COMMENT ON COLUMN team_members.role IS 'Роль пользователя в команде (MEMBER, LEAD, etc)';
COMMENT ON COLUMN team_members.created_at IS 'Время добавления пользователя в команду';

CREATE TABLE IF NOT EXISTS pull_requests (
    id UUID PRIMARY KEY,
    external_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    author_id UUID REFERENCES users(id) NOT NULL,
    status VARCHAR(20) DEFAULT 'OPEN',
    need_more_reviewers BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP NULL
);

COMMENT ON TABLE pull_requests IS 'Pull Requestы';
COMMENT ON COLUMN pull_requests.id IS 'Уникальный идентификатор Pull Request';
COMMENT ON COLUMN pull_requests.title IS 'Заголовок Pull Request';
COMMENT ON COLUMN pull_requests.description IS 'Описание Pull Request (опционально)';
COMMENT ON COLUMN pull_requests.author_id IS 'Идентификатор автора Pull Request';
COMMENT ON COLUMN pull_requests.status IS 'Статус PR: OPEN - открыт, MERGED - смержен';
COMMENT ON COLUMN pull_requests.need_more_reviewers IS 'Флаг необходимости дополнительных ревьюверов';
COMMENT ON COLUMN pull_requests.created_at IS 'Время создания Pull Request';
COMMENT ON COLUMN pull_requests.updated_at IS 'Время последнего обновления Pull Request';
COMMENT ON COLUMN pull_requests.merged_at IS 'Время мержа Pull Request (если смержен)';

CREATE TABLE IF NOT EXISTS pr_reviewers (
    id UUID PRIMARY KEY,
    pr_id UUID REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id UUID REFERENCES users(id) NOT NULL,
    team_id UUID REFERENCES teams(id) NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW(),
    replaced_at TIMESTAMP NULL,
    is_current BOOLEAN DEFAULT true,
    UNIQUE(pr_id, reviewer_id)
);

COMMENT ON TABLE pr_reviewers IS 'Назначенные ревьюверы на Pull Requestы';
COMMENT ON COLUMN pr_reviewers.id IS 'Уникальный идентификатор назначения';
COMMENT ON COLUMN pr_reviewers.pr_id IS 'Идентификатор Pull Request';
COMMENT ON COLUMN pr_reviewers.reviewer_id IS 'Идентификатор пользователя-ревьювера';
COMMENT ON COLUMN pr_reviewers.team_id IS 'Идентификатор команды, из которой назначен ревьювер';
COMMENT ON COLUMN pr_reviewers.assigned_at IS 'Время назначения ревьювера';
COMMENT ON COLUMN pr_reviewers.replaced_at IS 'Время замены ревьювера (если был заменен)';
COMMENT ON COLUMN pr_reviewers.is_current IS 'Флаг текущего назначения (false для замененных ревьюверов)';

CREATE TABLE IF NOT EXISTS pr_reviewer_history (
    id UUID PRIMARY KEY,
    pr_id UUID REFERENCES pull_requests(id) ON DELETE CASCADE,
    old_reviewer_id UUID REFERENCES users(id),
    new_reviewer_id UUID REFERENCES users(id),
    changed_by UUID REFERENCES users(id),
    changed_at TIMESTAMP DEFAULT NOW(),
    reason VARCHAR(100)
);

COMMENT ON TABLE pr_reviewer_history IS 'История изменений ревьюверов Pull Requestов';
COMMENT ON COLUMN pr_reviewer_history.id IS 'Уникальный идентификатор записи истории';
COMMENT ON COLUMN pr_reviewer_history.pr_id IS 'Идентификатор Pull Request';
COMMENT ON COLUMN pr_reviewer_history.old_reviewer_id IS 'Идентификатор старого ревьювера (NULL для первоначального назначения)';
COMMENT ON COLUMN pr_reviewer_history.new_reviewer_id IS 'Идентификатор нового ревьювера (NULL для снятия назначения)';
COMMENT ON COLUMN pr_reviewer_history.changed_by IS 'Идентификатор пользователя, выполнившего изменение';
COMMENT ON COLUMN pr_reviewer_history.changed_at IS 'Время изменения назначения';
COMMENT ON COLUMN pr_reviewer_history.reason IS 'Причина изменения: initial - первоначальное назначение, reassignment - переназначение, deactivation - деактивация пользователя';

CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_team_members_team_user ON team_members(team_id, user_id);
CREATE INDEX idx_team_members_user_team ON team_members(user_id, team_id);
CREATE INDEX idx_pull_requests_author ON pull_requests(author_id);
CREATE INDEX idx_pull_requests_merged_at ON pull_requests(merged_at) WHERE merged_at IS NOT NULL;
CREATE INDEX idx_pr_reviewers_pr_current ON pr_reviewers(pr_id, is_current);
CREATE INDEX idx_pr_reviewers_reviewer_current ON pr_reviewers(reviewer_id, is_current);
CREATE INDEX idx_pr_reviewers_team ON pr_reviewers(team_id);
CREATE INDEX idx_pr_reviewer_history_pr ON pr_reviewer_history(pr_id);
CREATE INDEX idx_pr_reviewer_history_reviewer ON pr_reviewer_history(old_reviewer_id, new_reviewer_id);
CREATE INDEX idx_pr_reviewer_history_changed_at ON pr_reviewer_history(changed_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewer_history;
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd