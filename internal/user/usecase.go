//go:generate mockgen -source usecase.go -destination mocks/usecase.go -package mocks
package user

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
)

type UserUsecase interface {
	Login(context.Context, string, string) (string, error)
	Register(context.Context, models.User) error
	UpdateUser(context.Context, models.User, models.User) error
	Withdraw(context.Context, int64) error
}
