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

type PRReviewerRepository struct {
	txMan txman.Manager
}

func NewPRReviewerRepository(txMan txman.Manager) *PRReviewerRepository {
	return &PRReviewerRepository{txMan: txMan}
}

func (r *PRReviewerRepository) GetPRReviewerByID(
	ctx context.Context,
	ID uuid.UUID,
) (data.PRReviewer, error) {
	var reviewer data.PRReviewer

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current
		FROM pr_reviewers
		WHERE id = $1
		`,
		ID,
	).Scan(
		&reviewer.ID,
		&reviewer.PullRequestID,
		&reviewer.ReviewerID,
		&reviewer.TeamID,
		&reviewer.AssignedAt,
		&reviewer.ReplacedAt,
		&reviewer.IsCurrent,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PRReviewer{}, errors.New(api.ErrNotFound)
		}
		return data.PRReviewer{}, errors.Wrap(err, errors.InternalError)
	}

	return reviewer, nil
}

func (r *PRReviewerRepository) CreatePRReviewer(
	ctx context.Context,
	reviewer data.PRReviewer,
) (data.PRReviewer, error) {
	var createdReviewer data.PRReviewer

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO pr_reviewers (id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current
		`,
		reviewer.ID,
		reviewer.PullRequestID,
		reviewer.ReviewerID,
		reviewer.TeamID,
		reviewer.AssignedAt,
		reviewer.ReplacedAt,
		reviewer.IsCurrent,
	).Scan(
		&createdReviewer.ID,
		&createdReviewer.PullRequestID,
		&createdReviewer.ReviewerID,
		&createdReviewer.TeamID,
		&createdReviewer.AssignedAt,
		&createdReviewer.ReplacedAt,
		&createdReviewer.IsCurrent,
	)
	if err != nil {
		return data.PRReviewer{}, errors.Wrap(err, errors.InternalError)
	}

	return createdReviewer, nil
}

func (r *PRReviewerRepository) GetCurrentReviewers(
	ctx context.Context,
	prID uuid.UUID,
) ([]data.PRReviewer, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current
		FROM pr_reviewers
		WHERE pr_id = $1 AND is_current = true
		`,
		prID,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalError)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var reviewers []data.PRReviewer
	for rows.Next() {
		var reviewer data.PRReviewer
		err := rows.Scan(
			&reviewer.ID,
			&reviewer.PullRequestID,
			&reviewer.ReviewerID,
			&reviewer.TeamID,
			&reviewer.AssignedAt,
			&reviewer.ReplacedAt,
			&reviewer.IsCurrent,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func (r *PRReviewerRepository) GetUserAssignedPRs(
	ctx context.Context,
	userID uuid.UUID,
) ([]data.PRReviewer, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current
		FROM pr_reviewers
		WHERE reviewer_id = $1 AND is_current = true
		`,
		userID,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalError)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var reviewers []data.PRReviewer
	for rows.Next() {
		var reviewer data.PRReviewer
		err := rows.Scan(
			&reviewer.ID,
			&reviewer.PullRequestID,
			&reviewer.ReviewerID,
			&reviewer.TeamID,
			&reviewer.AssignedAt,
			&reviewer.ReplacedAt,
			&reviewer.IsCurrent,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		reviewers = append(reviewers, reviewer)
	}

	return reviewers, nil
}

func (r *PRReviewerRepository) UpdatePRReviewer(
	ctx context.Context,
	reviewer data.PRReviewer,
) (data.PRReviewer, error) {
	var updatedReviewer data.PRReviewer

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE pr_reviewers 
		SET is_current = $1, replaced_at = $2
		WHERE id = $3
		RETURNING id, pr_id, reviewer_id, team_id, assigned_at, replaced_at, is_current
		`,
		reviewer.IsCurrent,
		reviewer.ReplacedAt,
		reviewer.ID,
	).Scan(
		&updatedReviewer.ID,
		&updatedReviewer.PullRequestID,
		&updatedReviewer.ReviewerID,
		&updatedReviewer.TeamID,
		&updatedReviewer.AssignedAt,
		&updatedReviewer.ReplacedAt,
		&updatedReviewer.IsCurrent,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PRReviewer{}, errors.New(api.ErrNotFound)
		}
		return data.PRReviewer{}, errors.Wrap(err, errors.InternalError)
	}

	return updatedReviewer, nil
}
