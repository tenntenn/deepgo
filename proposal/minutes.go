package proposal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

// ReviewMeetingMinute represents a single proposal review meeting minute.
type ReviewMeetingMinute struct {
	Body    string    `json:"body"`
	Created time.Time `json:"created,omitzero"`
	Updated time.Time `json:"updated,omitzero"`
}

// FetchReviewMeetingMinutesOptions defines parameters for fetching proposal review meeting minutes.
type FetchReviewMeetingMinutesOptions struct {
	Limit int
	Since time.Time
}

// IsLimitReached returns true if the given number n meets or exceeds the defined limit.
func (opts *FetchReviewMeetingMinutesOptions) IsLimitReached(n int) bool {
	return opts.Limit > 0 && n >= opts.Limit
}

// FetchReviewMeetingMinutes fetches proposal review meeting minutes from the golang/go issue tracker.
//
// It retrieves comments from issue #33502, which is used to record minutes of proposal review meetings.
// The results are limited by the provided options.
func FetchReviewMeetingMinutes(ctx context.Context, opts *FetchReviewMeetingMinutesOptions) ([]*ReviewMeetingMinute, error) {

	client := github.NewClient(nil)
	listopts := &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "desc",
		Since:     opts.Since,
	}

	var comments []*github.IssueComment
	for {
		page, resp, err := client.Issues.ListComments(ctx, "golang", "go", 33502, listopts)
		if err != nil {
			return nil, fmt.Errorf("failed to get issue comments: %w", err)
		}

		comments = append(comments, page...)

		if resp.NextPage == 0 {
			break
		}
		
		if opts.IsLimitReached(len(comments)) {
			comments = comments[:opts.Limit]
		}

		listopts.Page = resp.NextPage
	}

	minutes := make([]*ReviewMeetingMinute, len(comments))
	for i := range comments {
		minutes[i] = &ReviewMeetingMinute{
			Body:    comments[i].GetBody(),
			Created: comments[i].GetCreatedAt(),
			Updated: comments[i].GetUpdatedAt(),
		}
	}

	return minutes, nil
}
