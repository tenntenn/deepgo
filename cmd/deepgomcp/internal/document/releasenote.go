package document

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo/document"
)

func NewReleaseNoteTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("document/release-note",
			mcp.WithDescription("The tool gets the release note of given Go via https://go.dev/doc"),
			mcp.WithNumber("version",
				mcp.Description("the Go language version (e.g., go1.24)"),
			),
		),
		Handler: handleReleaseNote,
	}
}

func handleReleaseNote(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	langver, ok := request.Params.Arguments["version"].(string)
	if !ok {
		return mcp.NewToolResultError("the Go version (e.g, go1.24) must be specified."), nil
	}

	releaseNote, err := document.FetchReleaseNote(ctx, langver)
	if err != nil {
		return nil, fmt.Errorf("failed to get the release note: %w", err)
	}

	jsonRawMessage, err := json.Marshal(releaseNote)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonRawMessage)), nil
}
