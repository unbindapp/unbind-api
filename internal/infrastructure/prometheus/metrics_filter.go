package prometheus

import (
	"fmt"

	"github.com/google/uuid"
)

type MetricsFilter struct {
	TeamID        uuid.UUID
	ProjectID     uuid.UUID
	EnvironmentID uuid.UUID
	ServiceID     uuid.UUID
}

// buildLabelSelector constructs a Prometheus label selector string based on the provided MetricsFilter.
func buildLabelSelector(filter *MetricsFilter) string {
	selector := ""
	if filter.TeamID != uuid.Nil {
		selector += fmt.Sprintf("unbind_team=\"%s\",", filter.TeamID.String())
	}
	if filter.ProjectID != uuid.Nil {
		selector += fmt.Sprintf("unbind_project=\"%s\",", filter.ProjectID.String())
	}
	if filter.EnvironmentID != uuid.Nil {
		selector += fmt.Sprintf("unbind_environment=\"%s\",", filter.EnvironmentID.String())
	}
	if filter.ServiceID != uuid.Nil {
		selector += fmt.Sprintf("unbind_service=\"%s\",", filter.ServiceID.String())
	}

	// Remove trailing comma if present
	if len(selector) > 0 && selector[len(selector)-1] == ',' {
		selector = selector[:len(selector)-1]
	}

	if selector != "" {
		return "{" + selector + "}"
	}
	return ""
}
