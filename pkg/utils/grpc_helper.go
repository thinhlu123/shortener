package utils

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func GetFromMetadata(ctx context.Context, key string) string {
	var values []string

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values = md.Get(key)
	}

	if len(values) > 0 {
		return values[0]
	}

	return ""
}
