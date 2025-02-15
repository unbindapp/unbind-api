package kubeclient

// ? "Team" is synonymous with a Kubernetes namespace

import (
	"context"
	"fmt"
	"sort"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A team has a name and an underlying namespace
type UnbindTeam struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUnbindTeams returns a slice of UnbindTeam structs
func (k *KubeClient) GetUnbindTeams() ([]UnbindTeam, error) {
	namespaceRes := k.client.Resource(schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "namespaces",
	})

	namespaces, err := namespaceRes.List(context.Background(), metav1.ListOptions{
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
