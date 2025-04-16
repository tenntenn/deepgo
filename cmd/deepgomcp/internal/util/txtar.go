package util

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo/toolutil"
)

func NewCopyTxtarTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("util/copy-txtar",
			mcp.WithDescription("The tool copy files to given directory from txtar format string."),
			mcp.WithString("dir",
				mcp.Description("the dir parameter represents destination of files which must be absoluted path"),
				mcp.Required(),
			),
			mcp.WithString("txtar",
				mcp.Description("the txtar parameter represents txtar fomrat string"),
				mcp.Required(),
			),
		),
		Handler: handleCopyTxtar,
	}
}

func handleCopyTxtar(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dir, _ := request.Params.Arguments["dir"].(string)
	txtar, _ := request.Params.Arguments["txtar"].(string)

	if err := toolutil.CopyTxtar(dir, txtar); err != nil {
		return nil, fmt.Errorf("failed to copy files from txtar format string")
	}

	msg := fmt.Sprintf("copy files to %s", dir)
	return mcp.NewToolResultText(msg), nil
}
