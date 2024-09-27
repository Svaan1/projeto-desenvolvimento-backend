package db

import (
	"context"
)

type Database interface {
	NewDatabase(ctx context.Context) (Database, error)
	GetObject(ctx context.Context, key string, obj interface{}) error
	SetObject(ctx context.Context, key string, obj interface{}) error
}
