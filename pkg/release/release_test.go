package release

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v69/github"
	"github.com/stretchr/testify/suite"
)

// mockGitHubClient is a mock implementation of the GitHub client
type mockGitHubClient struct {
	listTagsFunc     func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error)
	listReleasesFunc func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error)
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

func (m *mockRepositoriesService) ListReleases(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
	return m.client.listReleasesFunc(ctx, owner, repo, opts)
}

type ReleaseTestSuite struct {
	suite.Suite
	manager  *Manager
	server   *httptest.Server
	metadata VersionMetadataMap
}

func (s *ReleaseTestSuite) SetupTest() {
	// Create test metadata
	s.metadata = VersionMetadataMap{
		"v0.0.1": {
			Version:     "v0.0.1",
			Description: "Initial release",
			Breaking:    false,
		},
		"v0.0.2": {
			Version:     "v0.0.2",
			Description: "Feature update",
			Breaking:    false,
		},
		"v0.0.3": {
			Version:     "v0.0.3",
			Description: "Bug fix",
			Breaking:    false,
		},
		"v0.1.0": {
			Version:     "v0.1.0",
			Description: "Major update",
			Breaking:    true,
			DependsOn:   []string{"v0.0.3"},
		},
	}

	// Create test server
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/unbindapp/unbind-releases/master/metadata.json" {
			data, _ := json.Marshal(s.metadata)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	// Create mock client that returns our test tags
	mockClient := &mockGitHubClient{
		listTagsFunc: func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
			return []*github.RepositoryTag{
				{Name: github.Ptr("v0.0.1")},
				{Name: github.Ptr("v0.0.2")},
				{Name: github.Ptr("v0.0.3")},
				{Name: github.Ptr("v0.1.0")},
			}, nil, nil
		},
		listReleasesFunc: func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
			return []*github.RepositoryRelease{
				{TagName: github.Ptr("v0.0.1")},
				{TagName: github.Ptr("v0.0.2")},
				{TagName: github.Ptr("v0.0.3")},
				{TagName: github.Ptr("v0.1.0")},
			}, nil, nil
		},
	}

	// Create manager with mock client and override the metadata URL
	s.manager = NewManager(mockClient, "unbindapp/unbind-releases")
	s.manager.metadataURL = s.server.URL + "/unbindapp/unbind-releases/master/metadata.json"
}

func (s *ReleaseTestSuite) TearDownTest() {
	s.server.Close()
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
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0"},
			expectError:    false,
		},
		{
			name:           "version without v prefix",
			currentVersion: "0.0.1",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0"},
			expectError:    false,
		},
		{
			name:           "latest version",
			currentVersion: "v0.1.0",
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
		name         string
		mockTags     []*github.RepositoryTag
		mockReleases []*github.RepositoryRelease
		expected     string
		expectError  bool
	}{
		{
			name: "valid tags with releases",
			mockTags: []*github.RepositoryTag{
				{Name: github.Ptr("v0.0.1")},
				{Name: github.Ptr("v0.0.2")},
				{Name: github.Ptr("v0.0.3")},
				{Name: github.Ptr("v0.1.0")},
			},
			mockReleases: []*github.RepositoryRelease{
				{TagName: github.Ptr("v0.0.1")},
				{TagName: github.Ptr("v0.0.2")},
				{TagName: github.Ptr("v0.0.3")},
				{TagName: github.Ptr("v0.1.0")},
			},
			expected:    "v0.1.0",
			expectError: false,
		},
		{
			name:         "no tags",
			mockTags:     []*github.RepositoryTag{},
			mockReleases: []*github.RepositoryRelease{},
			expected:     "",
			expectError:  true,
		},
		{
			name:         "only invalid tags",
			mockTags:     []*github.RepositoryTag{{Name: github.Ptr("invalid-tag")}},
			mockReleases: []*github.RepositoryRelease{},
			expected:     "",
			expectError:  true,
		},
		{
			name: "tags without releases",
			mockTags: []*github.RepositoryTag{
				{Name: github.Ptr("v0.0.1")},
				{Name: github.Ptr("v0.0.2")},
			},
			mockReleases: []*github.RepositoryRelease{},
			expected:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Update mock client for this test
			s.manager.client.(*mockGitHubClient).listTagsFunc = func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
				return tt.mockTags, nil, nil
			}
			s.manager.client.(*mockGitHubClient).listReleasesFunc = func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
				return tt.mockReleases, nil, nil
			}

			version, err := s.manager.GetLatestVersion(context.Background())
			if tt.expectError {
				s.Error(err)
				s.Empty(version)
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
			targetVersion:  "v0.0.3",
			expected:       []string{"v0.0.2", "v0.0.3"},
			expectError:    false,
		},
		{
			name:           "versions without v prefix",
			currentVersion: "0.0.1",
			targetVersion:  "0.0.3",
			expected:       []string{"v0.0.2", "v0.0.3"},
			expectError:    false,
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
			targetVersion:  "v0.1.0",
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
		{
			name:           "target version without release",
			currentVersion: "v0.0.1",
			targetVersion:  "v0.1.1",
			expected:       []string{},
			expectError:    false,
		},
		{
			name:           "path with dependencies",
			currentVersion: "v0.0.1",
			targetVersion:  "v0.1.0",
			expected:       []string{"v0.0.2", "v0.0.3", "v0.1.0"},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Reset mock client to default behavior for each test
			s.manager.client.(*mockGitHubClient).listTagsFunc = func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryTag, *github.Response, error) {
				return []*github.RepositoryTag{
					{Name: github.Ptr("v0.0.1")},
					{Name: github.Ptr("v0.0.2")},
					{Name: github.Ptr("v0.0.3")},
					{Name: github.Ptr("v0.1.0")},
				}, nil, nil
			}
			s.manager.client.(*mockGitHubClient).listReleasesFunc = func(ctx context.Context, owner, repo string, opts *github.ListOptions) ([]*github.RepositoryRelease, *github.Response, error) {
				return []*github.RepositoryRelease{
					{TagName: github.Ptr("v0.0.1")},
					{TagName: github.Ptr("v0.0.2")},
					{TagName: github.Ptr("v0.0.3")},
					{TagName: github.Ptr("v0.1.0")},
				}, nil, nil
			}

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
