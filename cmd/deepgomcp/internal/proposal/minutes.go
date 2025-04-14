package proposal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/tenntenn/deepgo/proposal"
)

func NewReviewMeetingMinutesTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("proposal/review-meeting-minutes",
			mcp.WithDescription("The tool gets the latest Go proposal weekly meeting minutes via the GitHub issue #33502"),
			mcp.WithNumber("limit",
				mcp.Description("the limit for the number of meeting minutes"),
			),
			mcp.WithString("since",
				mcp.Description("Filter meeting minutes from a specific date/time (YYYY-MM-DD hh:mm:ss)"),
			),
		),
		Handler: handleReviewMeetingMinutes,
	}
}

func handleReviewMeetingMinutes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var opts proposal.FetchReviewMeetingMinutesOptions
	limit, _ := request.Params.Arguments["limit"].(float64)
	opts.Limit = int(limit)
	sinceStr, ok := request.Params.Arguments["since"].(string)
	if sinceStr != "" || ok {
		since, err := time.Parse(time.DateTime, sinceStr)
		if err != nil {
			return mcp.NewToolResultError("invalid since format"), nil
		}
		opts.Since = since
	}
	minutes, err := proposal.FetchReviewMeetingMinutes(ctx, &opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get the proposal meeting minutes: %w", err)
	}

	jsonRawMessage, err := json.Marshal(minutes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return mcp.NewToolResultText(string(jsonRawMessage)), nil
}
