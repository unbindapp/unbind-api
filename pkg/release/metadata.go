package release

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/mod/semver"
)

// VersionMetadata represents the metadata for a specific version
type VersionMetadata struct {
	Version      string   `json:"version"`
	DependsOn    []string `json:"depends_on,omitempty"`    // Versions that must be installed before this one
	RequiredBy   []string `json:"required_by,omitempty"`   // Versions that require this version
	Breaking     bool     `json:"breaking,omitempty"`      // Whether this version contains breaking changes
	Description  string   `json:"description,omitempty"`   // Description of the version
	ReleaseNotes string   `json:"release_notes,omitempty"` // URL to release notes
}

// VersionMetadataMap maps version strings to their metadata
type VersionMetadataMap map[string]VersionMetadata

// GetVersionMetadata fetches the metadata for all versions
func (m *Manager) GetVersionMetadata(ctx context.Context) (VersionMetadataMap, error) {
	// Fetch the metadata file
	resp, err := http.Get(m.metadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch metadata: status code %d", resp.StatusCode)
	}

	// Read and parse the metadata
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata VersionMetadataMap
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return metadata, nil
}

// GetNextAvailableVersion returns the next version that can be updated to from the current version
func (m *Manager) GetNextAvailableVersion(ctx context.Context, currentVersion string) (string, error) {
	// Ensure current version has v prefix
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}

	// Get all available updates
	updates, err := m.AvailableUpdates(ctx, currentVersion)
	if err != nil {
		return "", err
	}

	// Get version metadata
	metadata, err := m.GetVersionMetadata(ctx)
	if err != nil {
		return "", err
	}

	// Find the first version that can be updated to
	for _, version := range updates {
		// Skip if version doesn't have metadata
		versionMeta, exists := metadata[version]
		if !exists {
			continue
		}

		// Check if this version has any dependencies
		if len(versionMeta.DependsOn) == 0 {
			return version, nil
		}

		// Check if all dependencies are satisfied
		canUpdate := true
		for _, dep := range versionMeta.DependsOn {
			// If the dependency is newer than current version, we can't update yet
			if semver.Compare(dep, currentVersion) > 0 {
				canUpdate = false
				break
			}
		}

		if canUpdate {
			return version, nil
		}

		// If we can't update to this version due to dependencies, return an error
		return "", fmt.Errorf("cannot update to version %s: requires version %s", version, versionMeta.DependsOn[0])
	}

	return "", fmt.Errorf("no available versions to update to")
}

// GetUpdatePath returns an ordered list of versions needed to update from current to target version
func (self *Manager) GetUpdatePath(ctx context.Context, currentVersion, targetVersion string) ([]string, error) {
	// Ensure versions have v prefix
	if !strings.HasPrefix(currentVersion, "v") {
		currentVersion = "v" + currentVersion
	}
	if !strings.HasPrefix(targetVersion, "v") {
		targetVersion = "v" + targetVersion
	}

	// Validate versions
	if !semver.IsValid(currentVersion) || !semver.IsValid(targetVersion) {
		return nil, fmt.Errorf("invalid version format")
	}

	// Get version metadata
	metadata, err := self.GetVersionMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get version metadata: %w", err)
	}

	// Get available updates
	updates, err := self.AvailableUpdates(ctx, currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get available updates: %w", err)
	}

	// Check if target version exists in updates
	targetExists := false
	for _, version := range updates {
		if version == targetVersion {
			targetExists = true
			break
		}
	}
	if !targetExists {
		return []string{}, nil
	}

	// Build the update path
	path := make([]string, 0)
	visited := make(map[string]bool)

	// First, add all versions between current and target
	for _, version := range updates {
		if semver.Compare(version, currentVersion) > 0 && semver.Compare(version, targetVersion) <= 0 {
			path = append(path, version)
			visited[version] = true
		}
	}

	// Sort the path by version
	sort.Slice(path, func(i, j int) bool {
		return semver.Compare(path[i], path[j]) < 0
	})

	// Now check dependencies and ensure they're in the path
	for i := 0; i < len(path); i++ {
		version := path[i]
		versionMeta, exists := metadata[version]
		if !exists {
			continue
		}

		// Check each dependency
		for _, dep := range versionMeta.DependsOn {
			// If dependency is not in path and is between current and target
			if !visited[dep] && semver.Compare(dep, currentVersion) > 0 && semver.Compare(dep, targetVersion) <= 0 {
				// Insert dependency before this version
				path = append(path[:i], append([]string{dep}, path[i:]...)...)
				visited[dep] = true
				i++ // Skip the dependency we just added
			}
		}
	}

	return path, nil
}
