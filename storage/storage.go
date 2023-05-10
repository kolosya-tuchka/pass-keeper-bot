package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, userService *UserService) error
	Pick(ctx context.Context, userID int64, service string) (*UserService, error)
	Remove(ctx context.Context, userID int64, service string) error
	Exists(ctx context.Context, userID int64, service string) (bool, error)
}

var (
	ErrNoSuchService = errors.New("no such service")
)

type UserService struct {
	UserID   int64
	Service  string
	Login    string
	Password string
}
