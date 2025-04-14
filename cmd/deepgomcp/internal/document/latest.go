package document

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tenntenn/goversion"
)

func NewLatestGoVersionTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("document/latest-go-version",
			mcp.WithDescription("The tool gets the latest note Go version via https://go.dev/VERSION?m=text"),
		),
		Handler: handleLatestGoVersion,
	}
}

func handleLatestGoVersion(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	latest, err := goversion.FetchLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest go version: %w", err)
	}

	jsonRawMessage, err := json.Marshal(latest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonRawMessage)), nil
}
