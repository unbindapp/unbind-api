package prometheus

import (
	"fmt"

	"github.com/google/uuid"
)

type MetricsFilterSumBy string

const (
	MetricsFilterSumByTeam        MetricsFilterSumBy = "team"
	MetricsFilterSumByProject     MetricsFilterSumBy = "project"
	MetricsFilterSumByEnvironment MetricsFilterSumBy = "environment"
	MetricsFilterSumByService     MetricsFilterSumBy = "service"
)

func (m MetricsFilterSumBy) Label() string {
	switch m {
	case MetricsFilterSumByTeam:
		return "label_unbind_team"
	case MetricsFilterSumByProject:
		return "label_unbind_project"
	case MetricsFilterSumByEnvironment:
		return "label_unbind_environment"
	case MetricsFilterSumByService:
		return "label_unbind_service"
	default:
		return "label_unbind_service"
	}
}

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
