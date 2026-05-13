package hub

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"smarthome/hub/core/logger"
	"smarthome/hub/pkg/updater"
	"time"
)

// UpdateConfig holds update configuration
type UpdateConfig struct {
	Enabled       bool
	AutoApply     bool
	Owner         string
	Repo          string
	GitHubToken   string
	CheckInterval time.Duration
}

// CheckAndApplyUpdates checks for updates and applies them if available
func CheckAndApplyUpdates(ctx context.Context, config UpdateConfig) error {
	if !config.Enabled {
		return nil
	}

	if config.Owner == "" || config.Repo == "" {
		logger.Log.Warn("Update check disabled: missing owner or repo configuration")
		return nil
	}

	logger.Log.Info("Checking for updates from %s/%s", config.Owner, config.Repo)

	source := updater.GitHubSource{
		Owner:    config.Owner,
		Repo:     config.Repo,
		APIToken: config.GitHubToken,
	}

	release, err := source.Latest(ctx)
	if err != nil {
		logger.Log.Warn("Failed to check for updates: %s", err.Error())
		return nil // Don't fail startup due to update check
	}

	if release == nil {
		logger.Log.Debug("No releases found")
		return nil
	}

	logger.Log.Info("Latest version available: %s", release.Version)

	if !config.AutoApply {
		logger.Log.Info("Auto-apply disabled. Manual update required.")
		return nil
	}

	// Apply updates
	exePath, err := os.Executable()
	if err != nil {
		logger.Log.Warn("Failed to get executable path: %s", err.Error())
		return nil
	}

	rootDir := filepath.Dir(exePath)

	client := updater.Client{
		Root:    rootDir,
		BaseURL: release.BaseURL,
		GOOS:    release.GOOS,
		GOARCH:  release.GOARCH,
	}

	changes, err := client.Plan(release.Manifest)
	if err != nil {
		logger.Log.Warn("Failed to plan updates: %s", err.Error())
		return nil
	}

	if len(changes) == 0 {
		logger.Log.Info("Already up to date")
		return nil
	}

	logger.Log.Info("Applying %d updates", len(changes))
	if err := client.Apply(ctx, release.Manifest); err != nil {
		logger.Log.Error("Failed to apply updates: %s", err.Error())
		return fmt.Errorf("update failed: %w", err)
	}

	logger.Log.Info("Updates applied successfully. Please restart the application.")
	return nil
}
