package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v69/github"
	"github.com/unbindapp/unbind-api/ent"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/internal/log"
	"github.com/unbindapp/unbind-api/internal/utils"
)

func (self *GithubClient) ReadUserAdminOrganizations(ctx context.Context, installation *ent.GithubInstallation) ([]*github.Organization, error) {
	if installation == nil || installation.AccountType != githubinstallation.AccountTypeUser || installation.Edges.GithubApp == nil {
		return nil, fmt.Errorf("Invalid installation")
	}

	privateKey, err := utils.DecodePrivateKey(installation.Edges.GithubApp.PrivateKey)
	if err != nil {
		return nil, err
	}

	bearerToken, err := utils.GenerateGithubJWT(installation.Edges.GithubApp.ID, privateKey)
	if err != nil {
		return nil, err
	}

	// Add token to client
	client := self.client.WithAuthToken(bearerToken)

	// Get user's organizations
	orgs, _, err := client.Organizations.List(ctx, installation.AccountLogin, nil)
	if err != nil {
		return nil, fmt.Errorf("Error getting user organizations: %v", err)
	}

	var adminOrgs []*github.Organization
	for _, org := range orgs {
		// Check membership status in the organization
		membership, _, err := client.Organizations.GetOrgMembership(ctx, installation.AccountLogin, org.GetLogin())
		if err != nil {
			// Log the error but continue processing other orgs
			log.Errorf("Error getting membership for org %s: %v\n", org.GetLogin(), err)
			continue
		}

		// Check if user has admin role in this organization
		if strings.EqualFold(membership.GetRole(), "admin") {
			adminOrgs = append(adminOrgs, org)
		}
	}
	return adminOrgs, nil
}
