package release

import (
	"context"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/suite"
)

// mockGitHubClient is a mock implementation of the GitHub client
type mockGitHubClient struct {
	listTagsFunc func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
}

func (m *mockGitHubClient) Repositories() RepositoriesServiceInterface {
	return &mockRepositoriesService{m}
}

type mockRepositoriesService struct {
	client *mockGitHubClient
}

func (m *mockRepositoriesService) ListTags(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
	return m.client.listTagsFunc(ctx, owner, repo, opts)
}

type ReleaseTestSuite struct {
	suite.Suite
	manager  *Manager
	mockTags []*github.RepositoryTag
}

func (s *ReleaseTestSuite) SetupTest() {
	// Create mock tags
	s.mockTags = []*github.RepositoryTag{
		{Name: github.String("v0.0.1")},
		{Name: github.String("v0.0.2")},
		{Name: github.String("v0.0.3")},
		{Name: github.String("v0.1.0")},
		{Name: github.String("v0.1.1")},
		{Name: github.String("v1.0.0")},
		{Name: github.String("invalid-tag")}, // Should be filtered out
	}

	// Create mock client
	mockClient := &mockGitHubClient{
		listTagsFunc: func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
			return s.mockTags, nil, nil
		},
	}

	// Create manager with mock client
	s.manager = NewManager(mockClient)
}

func (s *ReleaseTestSuite) TestAvailableUpdates() {
	tests := []struct {
		name           string
		currentVersion string
		expected       []string
		expectError    bool
	}{
		{
			name:           "valid version with updates",
			currentVersion: "v0.0.1",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0", "v0.1.1", "v1.0.0"},
			expectError:    false,
		},
		{
			name:           "version without v prefix",
			currentVersion: "0.0.1",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0", "v0.1.1", "v1.0.0"},
			expectError:    false,
		},
		{
			name:           "latest version",
			currentVersion: "v1.0.0",
			expected:       []string{},
			expectError:    false,
		},
		{
			name:           "invalid version",
			currentVersion: "invalid",
			expected:       []string{},
			expectError:    false,
		},
		{
			name:           "non-existent version",
			currentVersion: "v999.999.999",
			expected:       []string{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			updates, err := s.manager.AvailableUpdates(context.Background(), tt.currentVersion)
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, updates)
			}
		})
	}
}

func (s *ReleaseTestSuite) TestGetLatestVersion() {
	tests := []struct {
		name        string
		mockTags    []*github.RepositoryTag
		expected    string
		expectError bool
	}{
		{
			name:        "valid tags",
			mockTags:    s.mockTags,
			expected:    "v1.0.0",
			expectError: false,
		},
		{
			name:        "no tags",
			mockTags:    []*github.RepositoryTag{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "only invalid tags",
			mockTags:    []*github.RepositoryTag{{Name: github.String("invalid-tag")}},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Update mock client for this test
			s.manager.client.(*mockGitHubClient).listTagsFunc = func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
				return tt.mockTags, nil, nil
			}

			version, err := s.manager.GetLatestVersion(context.Background())
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, version)
			}
		})
	}
}

func (s *ReleaseTestSuite) TestGetUpdatePath() {
	tests := []struct {
		name           string
		currentVersion string
		targetVersion  string
		expected       []string
		expectError    bool
	}{
		{
			name:           "valid path",
			currentVersion: "v0.0.1",
			targetVersion:  "v0.1.0",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0"},
			expectError:    false,
		},
		{
			name:           "versions without v prefix",
			currentVersion: "0.0.1",
			targetVersion:  "0.1.0",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0"},
			expectError:    false,
		},
		{
			name:           "invalid current version",
			currentVersion: "invalid",
			targetVersion:  "v0.1.0",
			expected:       nil,
			expectError:    true,
		},
		{
			name:           "invalid target version",
			currentVersion: "v0.0.1",
			targetVersion:  "invalid",
			expected:       nil,
			expectError:    true,
		},
		{
			name:           "target version older than current",
			currentVersion: "v0.1.0",
			targetVersion:  "v0.0.1",
			expected:       []string{},
			expectError:    false,
		},
		{
			name:           "non-existent current version",
			currentVersion: "v999.999.999",
			targetVersion:  "v1.0.0",
			expected:       []string{},
			expectError:    false,
		},
		{
			name:           "non-existent target version",
			currentVersion: "v0.0.1",
			targetVersion:  "v999.999.999",
			expected:       []string{},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			path, err := s.manager.GetUpdatePath(context.Background(), tt.currentVersion, tt.targetVersion)
			if tt.expectError {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.expected, path)
			}
		})
	}
}

func TestReleaseSuite(t *testing.T) {
	suite.Run(t, new(ReleaseTestSuite))
}
