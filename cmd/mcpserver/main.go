package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/goproposal"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) (rerr error) {
	s := server.NewMCPServer(
		"Go Proposal MCP",
		"1.0.0",
		server.WithLogging(),
	)

	tool := mcp.NewTool("minutes",
		mcp.WithDescription("get the latest Go proposal meeting minutes via the GitHub issue #33502"),
		mcp.WithNumber("count",
			mcp.Required(),
			mcp.Description("the number of meeting minutes"),
		),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		count := request.Params.Arguments["count"].(float64)
		minutes, err := goproposal.FetchMinutes(ctx, int(count))
		if err != nil {
			return nil, fmt.Errorf("failed to get the proposal meeting minutes: %w", err)
		}

		jsonRawMessage, err := json.Marshal(minutes)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
		}

		return mcp.NewToolResultText(string(jsonRawMessage)), nil
	})

	slog.Info("run mcp server")
	if err := server.ServeStdio(s); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
