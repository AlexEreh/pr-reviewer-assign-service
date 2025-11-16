package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	goerrors "errors"

	"pr-reviewer-assign-service/internal/app/data"
	"pr-reviewer-assign-service/internal/app/delivery/http/api"
	"pr-reviewer-assign-service/pkg/errors"
	"pr-reviewer-assign-service/pkg/txman"
)

type PullRequestRepository struct {
	txMan txman.Manager
}

func NewPullRequestRepository(txMan txman.Manager) *PullRequestRepository {
	return &PullRequestRepository{txMan: txMan}
}

func (r *PullRequestRepository) GetPullRequestByID(
	ctx context.Context,
	ID uuid.UUID,
) (data.PullRequest, error) {
	var pr data.PullRequest

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		FROM pull_requests
		WHERE id = $1
		`,
		ID,
	).Scan(
		&pr.ID,
		&pr.ExternalID,
		&pr.Title,
		&pr.Description,
		&pr.AuthorID,
		&pr.Status,
		&pr.NeedMoreReviewers,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PullRequest{}, errors.New(api.ErrNotFound)
		}
		return data.PullRequest{}, errors.Wrap(err, errors.InternalError)
	}

	return pr, nil
}

func (r *PullRequestRepository) GetPullRequestByExternalID(
	ctx context.Context,
	externalID string,
) (data.PullRequest, error) {
	var pr data.PullRequest

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		FROM pull_requests
		WHERE external_id = $1
		`,
		externalID,
	).Scan(
		&pr.ID,
		&pr.ExternalID,
		&pr.Title,
		&pr.Description,
		&pr.AuthorID,
		&pr.Status,
		&pr.NeedMoreReviewers,
		&pr.CreatedAt,
		&pr.UpdatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PullRequest{}, errors.New(api.ErrNotFound)
		}
		return data.PullRequest{}, errors.Wrap(err, errors.InternalError)
	}

	return pr, nil
}

func (r *PullRequestRepository) CreatePullRequest(
	ctx context.Context,
	pr data.PullRequest,
) (data.PullRequest, error) {
	var createdPR data.PullRequest

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO pull_requests (id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		`,
		pr.ID,
		pr.ExternalID,
		pr.Title,
		pr.Description,
		pr.AuthorID,
		pr.Status,
		pr.NeedMoreReviewers,
		pr.CreatedAt,
		pr.UpdatedAt,
		pr.MergedAt,
	).Scan(
		&createdPR.ID,
		&createdPR.ExternalID,
		&createdPR.Title,
		&createdPR.Description,
		&createdPR.AuthorID,
		&createdPR.Status,
		&createdPR.NeedMoreReviewers,
		&createdPR.CreatedAt,
		&createdPR.UpdatedAt,
		&createdPR.MergedAt,
	)
	if err != nil {
		return data.PullRequest{}, errors.Wrap(err, errors.InternalError)
	}

	return createdPR, nil
}

// UpdatePullRequest обновляет данные PR и возвращает обновленную запись
func (r *PullRequestRepository) UpdatePullRequest(
	ctx context.Context,
	pr data.PullRequest,
) (data.PullRequest, error) {
	var updatedPR data.PullRequest

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE pull_requests 
		SET title = $1, description = $2, status = $3, need_more_reviewers = $4, updated_at = $5, merged_at = $6
		WHERE id = $7
		RETURNING id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		`,
		pr.Title,
		pr.Description,
		pr.Status,
		pr.NeedMoreReviewers,
		time.Now(),
		pr.MergedAt,
		pr.ID,
	).Scan(
		&updatedPR.ID,
		&updatedPR.ExternalID,
		&updatedPR.Title,
		&updatedPR.Description,
		&updatedPR.AuthorID,
		&updatedPR.Status,
		&updatedPR.NeedMoreReviewers,
		&updatedPR.CreatedAt,
		&updatedPR.UpdatedAt,
		&updatedPR.MergedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PullRequest{}, errors.New(api.ErrNotFound)
		}
		return data.PullRequest{}, errors.Wrap(err, errors.InternalError)
	}

	return updatedPR, nil
}

// MergePullRequest помечает PR как мерженный
func (r *PullRequestRepository) MergePullRequest(
	ctx context.Context,
	prID uuid.UUID,
) (data.PullRequest, error) {
	var mergedPR data.PullRequest

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE pull_requests 
		SET status = 'MERGED', updated_at = $1, merged_at = $2
		WHERE id = $3
		RETURNING id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		`,
		time.Now(),
		time.Now(),
		prID,
	).Scan(
		&mergedPR.ID,
		&mergedPR.ExternalID,
		&mergedPR.Title,
		&mergedPR.Description,
		&mergedPR.AuthorID,
		&mergedPR.Status,
		&mergedPR.NeedMoreReviewers,
		&mergedPR.CreatedAt,
		&mergedPR.UpdatedAt,
		&mergedPR.MergedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PullRequest{}, errors.New(api.ErrNotFound)
		}
		return data.PullRequest{}, errors.Wrap(err, errors.InternalError)
	}

	return mergedPR, nil
}

// GetOpenPullRequestsByAuthor возвращает открытые PR автора
func (r *PullRequestRepository) GetOpenPullRequestsByAuthor(
	ctx context.Context,
	authorID uuid.UUID,
) ([]data.PullRequest, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		FROM pull_requests
		WHERE author_id = $1 AND status = 'OPEN'
		`,
		authorID,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalError)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var prs []data.PullRequest
	for rows.Next() {
		var pr data.PullRequest
		err := rows.Scan(
			&pr.ID,
			&pr.ExternalID,
			&pr.Title,
			&pr.Description,
			&pr.AuthorID,
			&pr.Status,
			&pr.NeedMoreReviewers,
			&pr.CreatedAt,
			&pr.UpdatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		prs = append(prs, pr)
	}

	return prs, nil
}

func (r *PullRequestRepository) GetPullRequestsByStatus(
	ctx context.Context,
	status string,
) ([]data.PullRequest, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, external_id, title, description, author_id, status, need_more_reviewers, created_at, updated_at, merged_at
		FROM pull_requests
		WHERE status = $1
		`,
		status,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalError)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var prs []data.PullRequest
	for rows.Next() {
		var pr data.PullRequest
		err := rows.Scan(
			&pr.ID,
			&pr.ExternalID,
			&pr.Title,
			&pr.Description,
			&pr.AuthorID,
			&pr.Status,
			&pr.NeedMoreReviewers,
			&pr.CreatedAt,
			&pr.UpdatedAt,
			&pr.MergedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		prs = append(prs, pr)
	}

	return prs, nil
}
