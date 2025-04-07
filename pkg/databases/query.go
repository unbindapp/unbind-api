package databases

import (
	"context"
	"encoding/json"
	"fmt"
)

const DB_CATEGORY = "databases"

type DatabaseList struct {
	Databases []string `json:"databases"`
}

// ListDatabases lists all available databases
func (self *DatabaseProvider) ListDatabases(ctx context.Context, tagVersion string) (*DatabaseList, error) {
	// Base version URL
	baseURL := fmt.Sprintf(BaseDatabaseURL, tagVersion)

	// Fetch the index file that contains all categories and definitions
	indexURL := fmt.Sprintf("%s/index.json", baseURL)

	indexBytes, err := self.fetchURL(ctx, indexURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template index: %w", err)
	}

	// Parse the index
	var dbList DatabaseList
	if err := json.Unmarshal(indexBytes, &dbList); err != nil {
		return nil, fmt.Errorf("failed to parse template index: %w", err)
	}

	return &dbList, nil
}
