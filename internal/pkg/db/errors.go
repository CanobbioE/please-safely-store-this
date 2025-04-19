package db

import (
	"database/sql"
	"errors"
)

// IsNotFoundError returns true if sql.ErrNoRows is in the error's stack tree.
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
