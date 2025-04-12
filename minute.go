package goproposal

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/google/go-github/github"
)

type Minute struct {
	Body    string    `json:"body"`
	Created time.Time `json:"created,omitzero"`
	Updated time.Time `json:"updated,omitzero"`
}

func FetchMinutes(ctx context.Context, count int) ([]*Minute, error) {

	if count <= 0 {
		return nil, nil
	}

	client := github.NewClient(nil)
	var opts github.IssueListCommentsOptions
	var comments []*github.IssueComment

	for {
		page, resp, err := client.Issues.ListComments(ctx, "golang", "go", 33502, &opts)
		if err != nil {
			return nil, fmt.Errorf("failed to get issue comments: %w", err)
		}

		comments = append(comments, page...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	slices.Reverse(comments)
	comments = comments[:min(count, len(comments))]
	minutes := make([]*Minute, len(comments))
	for i := range comments {
		minutes[i] = &Minute{
			Body:    comments[i].GetBody(),
			Created: comments[i].GetCreatedAt(),
			Updated: comments[i].GetUpdatedAt(),
		}
	}

	return minutes, nil
}
