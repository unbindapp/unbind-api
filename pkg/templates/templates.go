package templates

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
)

var (
	BaseTemplateURL = "https://raw.githubusercontent.com/unbindapp/unbind-service-templates/refs/tags/%s"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

// UnbindTemplateProvider fetches templates from GitHub
type UnbindTemplateProvider struct {
	client *http.Client
}

func NewUnbindTemplateProvider() *UnbindTemplateProvider {
	return &UnbindTemplateProvider{
		client: http.DefaultClient,
	}
}

// fetchURL fetches a URL
func (self *UnbindTemplateProvider) fetchURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := self.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
