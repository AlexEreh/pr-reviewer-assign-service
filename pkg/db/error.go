package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	CodeUniqueViolation = "23505"
)

func IsUniqueConstraintViolationError(err error) bool {
	return checkPgxConstraintError(err)
}

func checkPgxConstraintError(err error) bool {
	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == CodeUniqueViolation {
		return true
	}

	return false
}
