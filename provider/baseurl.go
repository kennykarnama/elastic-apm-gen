package provider

import (
	"context"
	"os"
)

type BaseURL interface {
	Get(ctx context.Context) string
}

type BaseUrlFromEnv struct{}

func (b *BaseUrlFromEnv) Get() string {
	return os.Getenv("FAKE_BASE_URL")
}
