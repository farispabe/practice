package main

import (
	"context"
	"time"
)

type UserRepository interface {
	GetByTypeAndState(ctx context.Context, request GetUsersByTypeRequest) (users []User, err error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (response string, err error)
	Set(ctx context.Context, key string, data string, ttl time.Duration) (err error)
}

type Notifier interface {
	Notify(ctx context.Context, identifier string, message string) (err error)
}
