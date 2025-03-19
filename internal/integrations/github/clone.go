package github

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func (self *GithubClient) CloneRepository(ctx context.Context, appID, installationID int64, appPrivateKey string, repoURL string, refName string) (string, error) {
	bearerToken, err := self.GetInstallationToken(ctx, appID, installationID, appPrivateKey)
	if err != nil {
		return "", err
	}

	//Make a temporary directory to clone to
	tmpDir, err := os.MkdirTemp("", "unbind-api-clone")
	if err != nil {
		return "", err
	}

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &http.BasicAuth{
			Username: "x-access-token",
			Password: bearerToken,
		},
		ReferenceName: plumbing.ReferenceName(refName),
	})

	if err != nil {
		return "", fmt.Errorf("failed to clone repository: %v", err)
	}

	return tmpDir, nil
}
