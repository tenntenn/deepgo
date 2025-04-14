package style

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo/style"
)

func NewModernizeTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("style/modernize",
			mcp.WithDescription("The tool gets source code of modernize analyzer in gopls internal. All Go users must follow its rules."),
		),
		Handler: handleModernize,
	}
}

func handleModernize(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	modernize, err := style.FetchLatestModernize(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get source code of latest modernize: %w", err)
	}

	jsonRawMessage, err := json.Marshal(modernize)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonRawMessage)), nil
}
