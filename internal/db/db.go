package db

import (
	"context"
	"errors"
)

var (
	ErrNil = errors.New("no matching record found in redis database")
)

type Database interface {
	GetObject(ctx context.Context, key string, obj interface{}) error
	SetObject(ctx context.Context, key string, obj interface{}) error
}
