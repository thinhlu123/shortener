//go:generate mockgen -source usecase.go -destination mocks/usecase.go -package mocks
package link

import (
	"context"
	"github.com/thinhlu123/shortener/internal/models"
)

type LinkUsecases interface {
	CreateShortLink(context.Context, string) (string, error)
	GetLink(context.Context, string) (string, error)
	GetListLink(context.Context, models.Link) ([]models.Link, error)
}
