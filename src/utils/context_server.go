package utils

import (
	"context"
	"time"
)

func CreateContextServerWithTimeout() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return ctx
}