package persistence

import (
	"context"
	"time"
)

type SessionDB interface {
	Insert(ctx context.Context, userId, key string, timeout time.Time) (string, error)
	Get(ctx context.Context, sessionId string) (Session, error)
	UpdateTimeout(ctx context.Context, sessionId string, timeout time.Time) error
	Delete(sessionId string) error
}
