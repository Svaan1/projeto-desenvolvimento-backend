package db

import (
	"context"
)

type Database interface {
	NewDatabase(ctx context.Context) (Database, error)
}
