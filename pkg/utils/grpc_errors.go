package utils

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUserExists   = errors.New("user already exists")
	ErrPassword     = errors.New("password not match")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpireToken  = errors.New("expiry token")
)

// ParseGRPCErrStatusCode Parse error and get code
func ParseGRPCErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments), errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, context.Canceled), errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrUserExists):
		return codes.AlreadyExists
	case errors.Is(err, ErrPassword):
		return codes.InvalidArgument
	case errors.Is(err, ErrInvalidToken), errors.Is(err, ErrExpireToken):
		return codes.Unauthenticated
	}
	return codes.Internal
}

// MapGRPCErrCodeToHttpStatus Map GRPC errors codes to http status
func MapGRPCErrCodeToHttpStatus(code codes.Code) int {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.AlreadyExists:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
