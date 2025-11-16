package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	goerrors "errors"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/pkg/errors"
	"pr-reviewer-assign-service/pkg/txman"
)

type TeamMemberRepository struct {
	txMan txman.Manager
}

func NewTeamMemberRepository(txMan txman.Manager) TeamMemberRepository {
	return TeamMemberRepository{txMan: txMan}
}

func (r *TeamMemberRepository) GetTeamMemberByID(
	ctx context.Context,
	ID uuid.UUID,
) (data.TeamMember, error) {
	var teamMember data.TeamMember

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, team_id, user_id, role, created_at
		FROM team_members
		WHERE id = $1
		`,
		ID,
	).Scan(
		&teamMember.ID,
		&teamMember.TeamID,
		&teamMember.UserID,
		&teamMember.Role,
		&teamMember.CreatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.TeamMember{}, errors.New(api.ErrNotFound)
		}
		return data.TeamMember{}, errors.Wrap(err, errors.InternalError)
	}

	return teamMember, nil
}

func (r *TeamMemberRepository) GetTeamMemberByTeamAndUser(
	ctx context.Context,
	teamID, userID uuid.UUID,
) (data.TeamMember, error) {
	var teamMember data.TeamMember

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, team_id, user_id, role, created_at
		FROM team_members
		WHERE team_id = $1 AND user_id = $2
		`,
		teamID, userID,
	).Scan(
		&teamMember.ID,
		&teamMember.TeamID,
		&teamMember.UserID,
		&teamMember.Role,
		&teamMember.CreatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.TeamMember{}, errors.New(api.ErrNotFound)
		}
		return data.TeamMember{}, errors.Wrap(err, errors.InternalError)
	}

	return teamMember, nil
}

func (r *TeamMemberRepository) CreateTeamMember(
	ctx context.Context,
	teamMember data.TeamMember,
) (data.TeamMember, error) {
	var createdTeamMember data.TeamMember

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO team_members (id, team_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, team_id, user_id, role, created_at
		`,
		teamMember.ID,
		teamMember.TeamID,
		teamMember.UserID,
		teamMember.Role,
		teamMember.CreatedAt,
	).Scan(
		&createdTeamMember.ID,
		&createdTeamMember.TeamID,
		&createdTeamMember.UserID,
		&createdTeamMember.Role,
		&createdTeamMember.CreatedAt,
	)
	if err != nil {
		return data.TeamMember{}, errors.Wrap(err, errors.InternalError)
	}

	return createdTeamMember, nil
}

func (r *TeamMemberRepository) UpdateTeamMemberRole(
	ctx context.Context,
	teamMemberID uuid.UUID,
	role string,
) (data.TeamMember, error) {
	var updatedTeamMember data.TeamMember

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE team_members 
		SET role = $1
		WHERE id = $2
		RETURNING id, team_id, user_id, role, created_at
		`,
		role,
		teamMemberID,
	).Scan(
		&updatedTeamMember.ID,
		&updatedTeamMember.TeamID,
		&updatedTeamMember.UserID,
		&updatedTeamMember.Role,
		&updatedTeamMember.CreatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.TeamMember{}, errors.New(api.ErrNotFound)
		}
		return data.TeamMember{}, errors.Wrap(err, errors.InternalError)
	}

	return updatedTeamMember, nil
}

func (r *TeamMemberRepository) DeleteTeamMember(ctx context.Context, teamMemberID uuid.UUID) error {
	result, err := r.txMan.Executor(ctx).ExecContext(
		ctx,
		`DELETE FROM team_members WHERE id = $1`,
		teamMemberID,
	)
	if err != nil {
		return errors.Wrap(err, errors.InternalError)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, errors.InternalError)
	}

	if rowsAffected == 0 {
		return errors.New(api.ErrNotFound)
	}

	return nil
}

func (r *TeamMemberRepository) GetTeamMembersByTeamID(
	ctx context.Context,
	teamID uuid.UUID,
) ([]data.TeamMember, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, team_id, user_id, role, created_at
		FROM team_members
		WHERE team_id = $1
		`,
		teamID,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return []data.TeamMember{}, errors.New(api.ErrNotFound)
		}
		return nil, errors.Wrap(err, errors.InternalError)
	}
	defer rows.Close()

	var teamMembers []data.TeamMember
	for rows.Next() {
		var teamMember data.TeamMember
		err := rows.Scan(
			&teamMember.ID,
			&teamMember.TeamID,
			&teamMember.UserID,
			&teamMember.Role,
			&teamMember.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		teamMembers = append(teamMembers, teamMember)
	}

	return teamMembers, nil
}

func (r *TeamMemberRepository) GetTeamMembersByUserID(
	ctx context.Context,
	teamID uuid.UUID,
) ([]data.TeamMember, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, team_id, user_id, role, created_at
		FROM team_members
		WHERE user_id = $1
		`,
		teamID,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return []data.TeamMember{}, errors.New(api.ErrNotFound)
		}
		return nil, errors.Wrap(err, errors.InternalError)
	}
	defer rows.Close()

	var teamMembers []data.TeamMember
	for rows.Next() {
		var teamMember data.TeamMember
		err := rows.Scan(
			&teamMember.ID,
			&teamMember.TeamID,
			&teamMember.UserID,
			&teamMember.Role,
			&teamMember.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		teamMembers = append(teamMembers, teamMember)
	}

	return teamMembers, nil
}
