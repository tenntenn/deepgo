package style

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo/style"
)

func NewGoStyleTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("style/gostyle",
			mcp.WithDescription("The tool gets the Go Style Guide. The Go Style Guide and accompanying documents codify the current best approaches for writing readable and idiomatic Go. Adherence to the Style Guide is not intended to be absolute, and these documents will never be exhaustive. Our intention is to minimize the guesswork of writing readable Go so that newcomers to the language can avoid common mistakes. The Style Guide also serves to unify the style guidance given by anyone reviewing Go code at Google."),
		),
		Handler: handleGoStyle,
	}
}

func handleGoStyle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gostyle, err := style.FetchGoStyle(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get source code of latest modernize: %w", err)
	}

	jsonRawMessage, err := json.Marshal(gostyle)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonRawMessage)), nil
}
