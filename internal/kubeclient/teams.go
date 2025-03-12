package kubeclient

// ? "Team" is synonymous with a Kubernetes namespace

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/unbindapp/unbind-api/internal/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A team has a name and an underlying namespace
type UnbindTeam struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUnbindTeams returns a slice of UnbindTeam structs
func (k *KubeClient) GetUnbindTeams(ctx context.Context, bearerToken string) ([]UnbindTeam, error) {
	client, err := k.createClientWithToken(bearerToken)
	if err != nil {
		log.Errorf("Error creating client with token: %v", err)
		return nil, fmt.Errorf("error creating client with token: %v", err)
	}

	// List namespaces with the unbind-team label
	namespaces, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: "unbind-team",
	})
	if err != nil {
		return nil, fmt.Errorf("error listing namespaces: %v", err)
	}

	teams := make([]UnbindTeam, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		teamValue := ns.GetLabels()["unbind-team"]
		if teamValue != "" {
			teams = append(teams, UnbindTeam{
				Name:      teamValue,
				Namespace: ns.GetName(),
				CreatedAt: ns.GetCreationTimestamp().Time,
			})
		}
	}

	// Sort by data created by default
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].CreatedAt.After(teams[j].CreatedAt)
	})

	return teams, nil
}
