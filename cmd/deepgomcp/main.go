package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tenntenn/deepgo/cmd/deepgomcp/internal"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	s, err := internal.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to create a MCP server: %w", err)
	}

	if err := s.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
