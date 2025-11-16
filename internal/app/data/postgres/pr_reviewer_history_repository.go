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

type PRReviewerHistoryRepository struct {
	txMan txman.Manager
}

func NewPRReviewerHistoryRepository(txMan txman.Manager) *PRReviewerHistoryRepository {
	return &PRReviewerHistoryRepository{txMan: txMan}
}

// GetPRReviewerHistoryByID возвращает запись истории по ID
func (r *PRReviewerHistoryRepository) GetPRReviewerHistoryByID(
	ctx context.Context,
	ID uuid.UUID,
) (data.PRReviewerHistory, error) {
	var history data.PRReviewerHistory

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, pr_id, old_reviewer_id, new_reviewer_id, changed_by, changed_at, reason
		FROM pr_reviewer_history
		WHERE id = $1
		`,
		ID,
	).Scan(
		&history.ID,
		&history.PullRequestID,
		&history.OldReviewerID,
		&history.NewReviewerID,
		&history.ChangedBy,
		&history.ChangedAt,
		&history.Reason,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.PRReviewerHistory{}, errors.New(api.ErrNotFound)
		}
		return data.PRReviewerHistory{}, errors.Wrap(err, errors.InternalError)
	}

	return history, nil
}

// CreatePRReviewerHistory создает новую запись истории и возвращает ее с ID
func (r *PRReviewerHistoryRepository) CreatePRReviewerHistory(
	ctx context.Context,
	history data.PRReviewerHistory,
) (data.PRReviewerHistory, error) {
	var createdHistory data.PRReviewerHistory

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO pr_reviewer_history (id, pr_id, old_reviewer_id, new_reviewer_id, changed_by, changed_at, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, pr_id, old_reviewer_id, new_reviewer_id, changed_by, changed_at, reason
		`,
		history.ID,
		history.PullRequestID,
		history.OldReviewerID,
		history.NewReviewerID,
		history.ChangedBy,
		history.ChangedAt,
		history.Reason,
	).Scan(
		&createdHistory.ID,
		&createdHistory.PullRequestID,
		&createdHistory.OldReviewerID,
		&createdHistory.NewReviewerID,
		&createdHistory.ChangedBy,
		&createdHistory.ChangedAt,
		&createdHistory.Reason,
	)
	if err != nil {
		return data.PRReviewerHistory{}, errors.Wrap(err, errors.InternalError)
	}

	return createdHistory, nil
}

// GetPRReviewerHistory возвращает историю изменений для PR
func (r *PRReviewerHistoryRepository) GetPRReviewerHistory(
	ctx context.Context,
	prID uuid.UUID,
) ([]data.PRReviewerHistory, error) {
	rows, err := r.txMan.Executor(ctx).QueryContext(
		ctx,
		`
		SELECT id, pr_id, old_reviewer_id, new_reviewer_id, changed_by, changed_at, reason
		FROM pr_reviewer_history
		WHERE pr_id = $1
		ORDER BY changed_at DESC
		`,
		prID,
	)
	if err != nil {
		return nil, errors.Wrap(err, errors.InternalError)
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var history []data.PRReviewerHistory
	for rows.Next() {
		var entry data.PRReviewerHistory
		err := rows.Scan(
			&entry.ID,
			&entry.PullRequestID,
			&entry.OldReviewerID,
			&entry.NewReviewerID,
			&entry.ChangedBy,
			&entry.ChangedAt,
			&entry.Reason,
		)
		if err != nil {
			return nil, errors.Wrap(err, errors.InternalError)
		}
		history = append(history, entry)
	}

	return history, nil
}
