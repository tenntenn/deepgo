package document

import (
	"context"
	"fmt"
	"go/version"
	"io"
	"net/http"

	"github.com/tenntenn/goversion"
)

// ReleaseNote represents a release note.
type ReleaseNote struct {
	Version string `json:"version"`
	Body    string `json:"body"`
}

// FetchReviewMeetingMinutes fetches proposal review meeting minutes from the golang/go issue tracker.
func FetchReleaseNote(ctx context.Context, gover string) (*ReleaseNote, error) {
	langver := version.Lang(gover)
	if langver == "" {
		return nil, fmt.Errorf("invalid go version: %q", gover)
	}

	latest, err := goversion.FetchLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest go version: %w", err)
	}

	if version.Compare(langver, version.Lang(latest.Version)) > 1 {
		return nil, fmt.Errorf("%q has not released yet", langver)
	}

	docURL := fmt.Sprintf("https://go.dev/doc/%s", langver)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, docURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request for %q: %w", docURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send a request to %q: %w", docURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid HTTP Status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &ReleaseNote{
		Version: langver,
		Body:    string(body),
	}, nil
}
