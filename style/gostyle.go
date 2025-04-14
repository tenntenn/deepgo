package style

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const gostyleBaseURL = "https://google.github.io/styleguide/go"

// GoStyle represents Go style guide in Google.
type GoStyle struct {
	Overview      string `json:"overview"`
	Guide         string `json:"guide"`
	Decisions     string `json:"decisions"`
	// Too long
	//BestPractices string `json:"best_practices"`
}

// FetchGoStyle fetches Go style guide in Google.
func FetchGoStyle(ctx context.Context) (*GoStyle, error) {

	overview, err := fetchGoStyleHTML(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get overview: %w", err)
	}

	guide, err := fetchGoStyleHTML(ctx, "/guide")
	if err != nil {
		return nil, fmt.Errorf("failed to get guide: %w", err)
	}

	decisions, err := fetchGoStyleHTML(ctx, "/decisions")
	if err != nil {
		return nil, fmt.Errorf("failed to get decisions: %w", err)
	}

	// Too long
	//bestPractices, err := fetchGoStyleHTML(ctx, "/best-practices")
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get best practices: %w", err)
	//}

	return &GoStyle{
		Overview:      overview,
		Guide:         guide,
		Decisions:     decisions,
		//BestPractices: bestPractices,
	}, nil
}

func fetchGoStyleHTML(ctx context.Context, sub string) (string, error) {

	url := gostyleBaseURL + "/" + sub
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request for %q: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send a request to %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid HTTP Status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}
