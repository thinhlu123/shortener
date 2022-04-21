//go:generate mockgen -source repository.go -destination mocks/repository.go -package mocks
package user

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
)

type UserRepository interface {
	Login(context.Context, models.User) (string, error)
	Register(context.Context, models.User) error
	UpdateUser(context.Context, models.User, models.User) error
}
