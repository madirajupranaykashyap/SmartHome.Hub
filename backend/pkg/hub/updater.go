package hub

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"smarthome/hub/pkg/updater"
	"strings"
)

// CheckForUpdates checks GitHub for new releases
func CheckForUpdates(ctx context.Context, owner, repo, token string) (updater.GitHubRelease, error) {
	source := updater.GitHubSource{
		Owner:    owner,
		Repo:     repo,
		APIToken: token,
	}

	release, err := source.Latest(ctx)
	if err != nil {
		return updater.GitHubRelease{}, err
	}

	return release, nil
}

// ApplyUpdates applies the update manifest if updates are available
func ApplyUpdates(ctx context.Context, release updater.GitHubRelease) ([]updater.Change, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	rootDir := filepath.Dir(exePath)

	// Construct base URL from manifest URL by removing the filename
	baseURL := release.ManifestURL
	lastSlash := strings.LastIndex(baseURL, "/")
	if lastSlash > 0 {
		baseURL = baseURL[:lastSlash+1]
	}

	client := updater.Client{
		Root:    rootDir,
		BaseURL: baseURL,
		GOOS:    runtime.GOOS,
		GOARCH:  runtime.GOARCH,
	}

	changes, err := client.Plan(release.Manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to plan updates: %w", err)
	}

	if len(changes) == 0 {
		return changes, nil
	}

	changes, err = client.Apply(ctx, release.Manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to apply updates: %w", err)
	}

	return changes, nil
}
