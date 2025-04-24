package github

import (
	"context"
	"fmt"
	"os"

	charmLog "github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

func (self *GithubClient) CloneRepository(ctx context.Context, appID, installationID int64, appPrivateKey string, repoURL string, refName string, commitSHA string) (string, error) {
	bearerToken, err := self.GetInstallationToken(ctx, appID, installationID, appPrivateKey)
	if err != nil {
		return "", err
	}

	// Make a temporary directory to clone to
	tmpDir, err := os.MkdirTemp("", "unbind-api-clone")
	if err != nil {
		return "", err
	}

	// Set up the logger
	logger := &loggerOutput{
		logger: log.GetLogger(),
	}

	cloneOptions := &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: bearerToken,
		},
		Progress: logger,
	}

	// If a specific commit SHA is provided, we need to modify the clone strategy
	if commitSHA != "" {
		// First clone without specifying a branch
		cloneOptions.ReferenceName = ""
		cloneOptions.SingleBranch = false
	} else {
		// Standard branch-based clone
		cloneOptions.SingleBranch = true
		cloneOptions.ReferenceName = plumbing.ReferenceName(refName)
	}

	repo, err := git.PlainClone(tmpDir, false, cloneOptions)
	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %v", err)
	}

	// If a specific commit SHA is provided, checkout that commit
	if commitSHA != "" {
		worktree, err := repo.Worktree()
		if err != nil {
			return "", fmt.Errorf("failed to get worktree: %v", err)
		}

		err = worktree.Checkout(&git.CheckoutOptions{
			Hash: plumbing.NewHash(commitSHA),
		})
		if err != nil {
			return "", fmt.Errorf("failed to checkout commit %s: %v", commitSHA, err)
		}
	}

	return tmpDir, nil
}

type loggerOutput struct {
	logger *charmLog.Logger
}

func (l *loggerOutput) Write(p []byte) (n int, err error) {
	l.logger.Infof("%s", p)
	return len(p), nil
}
