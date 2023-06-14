package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	cmd "go.szostok.io/codeowners/cmd/codeowners"
)

func main() {
	ctx, cancelFunc := WithStopContext(context.Background())
	defer cancelFunc()

	if err := cmd.RootCmd().ExecuteContext(ctx); err != nil {
		// error is already handled by `cobra`, we don't want to log it here as we will duplicate the message.
		// If needed, based on error type we can exit with different codes.
		//nolint:gocritic
		os.Exit(1)
	}
}

// WithStopContext returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed on of SIGINT or SIGTERM signals.
func WithStopContext(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
		case <-sigCh:
			cancel()
		}
	}()

	return ctx, cancel
}
