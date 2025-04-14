package style

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

const modernizeBaseURL = "https://proxy.golang.org/golang.org/x/tools/gopls"

// Modernize represents source codes of modernize.
type Modernize struct {
	Version string `json:"version"`
	Source  string `json:"source"`
}

// FetchLatestModernize fetches source code of modernize.Analyzer in gopls internal.
func FetchLatestModernize(ctx context.Context) (*Modernize, error) {

	latestVersion, err := latestModernizeVersion(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version of modernize: %w", err)
	}

	fsys, err := fetchSourceZIP(ctx, latestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get zip file of modernize@%s: %w", latestVersion, err)
	}

	ar, err := txtarfs.From(fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to txtar format: %w", err)
	}

	return &Modernize{
		Version: latestVersion,
		Source:  string(txtar.Format(ar)),
	}, nil
}

func latestModernizeVersion(ctx context.Context) (string, error) {

	latestURL := modernizeBaseURL + "/@latest"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request for %q: %w", latestURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send a request to %q: %w", latestURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid HTTP Status: %d", resp.StatusCode)
	}

	var info struct{ Version string }
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	return info.Version, nil
}

func fetchSourceZIP(ctx context.Context, latestVersion string) (fs.FS, error) {
	zipURL := fmt.Sprintf("%s/@v/%s.zip", modernizeBaseURL, latestVersion)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, zipURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request for %q: %w", zipURL, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send a request to %q: %w", zipURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid HTTP Status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create ZIP reader: %w", err)
	}

	fsys, err := fs.Sub(zr, "golang.org/x/tools/gopls@v0.18.1/internal/analysis/modernize")
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectory of modernize: %w", err)
	}

	return fsys, nil
}
