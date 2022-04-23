//go:generate mockgen -source repository.go -destination mocks/repository.go -package mocks
package link

import (
	"context"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/thinhlu123/shortener/internal/models"
)

type LinkRepository interface {
	CreateShortLink(context.Context, models.Link) error
	GetLink(context.Context, models.Link) (string, error)
	IncreaseClickCount(context.Context, models.Link) error
	GetListLink(context.Context, models.Link) ([]models.Link, error)
}
