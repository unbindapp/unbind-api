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

// buildLabelSelector constructs the selector for kube_pod_labels
func buildLabelSelector(filter *MetricsFilter) string {
	if filter == nil {
		return ""
	}

	var labelFilters []string

	if filter.TeamID != uuid.Nil {
		labelFilters = append(labelFilters, fmt.Sprintf(`label_unbind_team="%s"`, filter.TeamID.String()))
	}

	if filter.ProjectID != uuid.Nil {
		labelFilters = append(labelFilters, fmt.Sprintf(`label_unbind_project="%s"`, filter.ProjectID.String()))
	}

	if filter.EnvironmentID != uuid.Nil {
		labelFilters = append(labelFilters, fmt.Sprintf(`label_unbind_environment="%s"`, filter.EnvironmentID.String()))
	}

	if filter.ServiceID != uuid.Nil {
		labelFilters = append(labelFilters, fmt.Sprintf(`label_unbind_service="%s"`, filter.ServiceID.String()))
	}

	if len(labelFilters) == 0 {
		return ""
	}

	// Combine all filters with logical AND (comma in PromQL)
	filterQuery := "{"
	for i, labelFilter := range labelFilters {
		if i > 0 {
			filterQuery += ", "
		}
		filterQuery += labelFilter
	}
	filterQuery += "}"

	return filterQuery
}
