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

type UserRepository struct {
	txMan txman.Manager
}

func (r *UserRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (data.User, error) {
	var user data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, username, email, is_active, created_at, updated_at 
		FROM users
		WHERE id = $1
		`,
		ID,
	).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return user, nil
}

func (r *UserRepository) GetUserByExternalID(
	ctx context.Context,
	externalID string,
) (data.User, error) {
	var user data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, username, email, is_active, created_at, updated_at 
		FROM users
		WHERE external_id = $1
		`,
		externalID,
	).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return user, nil
}

func (r *UserRepository) GetUserByName(ctx context.Context, userName string) (data.User, error) {
	var user data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, username, email, is_active, created_at, updated_at 
		FROM users
		WHERE username = $1
		`,
		userName,
	).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (data.User, error) {
	var user data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		SELECT id, external_id, username, email, is_active, created_at, updated_at 
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user data.User) (data.User, error) {
	var createdUser data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		INSERT INTO users (id, external_id, username, email, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, external_id, username, email, is_active, created_at, updated_at
		`,
		user.ID,
		user.ExternalID,
		user.Username,
		user.Email,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(
		&createdUser.ID,
		&createdUser.ExternalID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.IsActive,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return createdUser, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user data.User) (data.User, error) {
	var updatedUser data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE users 
		SET username = $1, email = $2, is_active = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, external_id, username, email, is_active, created_at, updated_at
		`,
		user.Username,
		user.Email,
		user.IsActive,
		time.Now(),
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.ExternalID,
		&updatedUser.Username,
		&updatedUser.Email,
		&updatedUser.IsActive,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return updatedUser, nil
}

func (r *UserRepository) SetUserActive(
	ctx context.Context,
	userID uuid.UUID,
	isActive bool,
) (data.User, error) {
	var user data.User

	err := r.txMan.Executor(ctx).QueryRowContext(
		ctx,
		`
		UPDATE users 
		SET is_active = $1, updated_at = $2
		WHERE id = $3
		RETURNING id, external_id, username, email, is_active, created_at, updated_at
		`,
		isActive,
		time.Now(),
		userID,
	).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if goerrors.Is(err, sql.ErrNoRows) {
			return data.User{}, errors.New(api.ErrNotFound)
		}
		return data.User{}, errors.Wrap(err, errors.InternalError)
	}

	return user, nil
}
