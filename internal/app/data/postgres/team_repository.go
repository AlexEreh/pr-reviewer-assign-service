package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	goerrors "errors"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/pkg/db"
	"pr-reviewer-assign-service/pkg/errors"
	"pr-reviewer-assign-service/pkg/log"
	"pr-reviewer-assign-service/pkg/txman"
)

type TeamRepository struct {
	txMan txman.Manager
}

func NewTeamRepository(txMan txman.Manager) *TeamRepository {
	return &TeamRepository{txMan: txMan}
}

// GetTeamByID возвращает команду по внутреннему ID
func (r *TeamRepository) GetTeamByID(ctx context.Context, ID uuid.UUID) (data.Team, error) {
	var team data.Team

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, name, description, created_at, updated_at 
		FROM teams
		WHERE id = $1
		`,
		ID,
	).Scan(
		&team.ID,
		&team.ExternalID,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.Team{}, errors.New(api.ErrNotFound)
		}

		return data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	return team, nil
}

// GetTeamByExternalID возвращает команду по внешнему ID
func (r *TeamRepository) GetTeamByExternalID(
	ctx context.Context,
	externalID string,
) (data.Team, error) {
	var team data.Team

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, name, description, created_at, updated_at 
		FROM teams
		WHERE external_id = $1
		`,
		externalID,
	).Scan(
		&team.ID,
		&team.ExternalID,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.Team{}, errors.New(api.ErrNotFound)
		}

		return data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	return team, nil
}

// GetTeamByName возвращает команду по имени
func (r *TeamRepository) GetTeamByName(ctx context.Context, name string) (data.Team, error) {
	var team data.Team

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, name, description, created_at, updated_at 
		FROM teams
		WHERE name = $1
		`,
		name,
	).Scan(
		&team.ID,
		&team.ExternalID,
		&team.Name,
		&team.Description,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.Team{}, errors.New(api.ErrNotFound)
		}

		return data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	return team, nil
}

func (r *TeamRepository) GetAllTeams(ctx context.Context) ([]data.Team, error) {
	var teams []data.Team

	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, external_id, name, description, created_at, updated_at 
		FROM teams
		`,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return []data.Team{}, errors.New(api.ErrNotFound)
		}

		return []data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if rows.Err() != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return []data.Team{}, errors.New(api.ErrNotFound)
		}

		return []data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	var team data.Team
	for rows.Next() {
		err = rows.Scan(
			&team.ID,
			&team.ExternalID,
			&team.Name,
			&team.Description,
			&team.CreatedAt,
			&team.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// CreateTeam создает новую команду и возвращает ее с ID
func (r *TeamRepository) CreateTeam(ctx context.Context, team data.Team) (data.Team, error) {
	var createdTeam data.Team

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO teams (id, external_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, external_id, name, description, created_at, updated_at
		`,
		team.ID,
		team.ExternalID,
		team.Name,
		team.Description,
		team.CreatedAt,
		team.UpdatedAt,
	).Scan(
		&createdTeam.ID,
		&createdTeam.ExternalID,
		&createdTeam.Name,
		&createdTeam.Description,
		&createdTeam.CreatedAt,
		&createdTeam.UpdatedAt,
	)
	if err != nil {
		log.LoggerFromCtx(ctx).Error("err", zap.Error(err))
		if db.IsUniqueConstraintViolationError(err) {
			return data.Team{}, errors.New(api.ErrTeamExists)
		}

		return data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	return createdTeam, nil
}

// UpdateTeam обновляет данные команды и возвращает обновленную запись
func (r *TeamRepository) UpdateTeam(ctx context.Context, team data.Team) (data.Team, error) {
	var updatedTeam data.Team

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE teams 
		SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, external_id, name, description, created_at, updated_at
		`,
		team.Name,
		team.Description,
		time.Now(),
		team.ID,
	).Scan(
		&updatedTeam.ID,
		&updatedTeam.ExternalID,
		&updatedTeam.Name,
		&updatedTeam.Description,
		&updatedTeam.CreatedAt,
		&updatedTeam.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.Team{}, errors.New(api.ErrNotFound)
		}
		return data.Team{}, errors.Wrap(err, errors.InternalError)
	}

	return updatedTeam, nil
}
