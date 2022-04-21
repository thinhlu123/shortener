//go:generate mockgen -source usecase.go -destination mocks/usecase.go -package mocks
package link

import (
	"context"
)

type LinkUsecases interface {
	CreateShortLink(context.Context, string) (string, error)
	GetLink(context.Context, string) (string, error)
}
