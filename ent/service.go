// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/unbindapp/unbind-api/ent/deployment"
	"github.com/unbindapp/unbind-api/ent/environment"
	"github.com/unbindapp/unbind-api/ent/githubinstallation"
	"github.com/unbindapp/unbind-api/ent/service"
	"github.com/unbindapp/unbind-api/ent/serviceconfig"
)

// Service is the model entity for the Service schema.
type Service struct {
	config `json:"-"`
	// ID of the ent.
	// The primary key of the entity.
	ID uuid.UUID `json:"id"`
	// The time at which the entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// The time at which the entity was last updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// DisplayName holds the value of the "display_name" field.
	DisplayName string `json:"display_name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// EnvironmentID holds the value of the "environment_id" field.
	EnvironmentID uuid.UUID `json:"environment_id,omitempty"`
	// Optional reference to GitHub installation
	GithubInstallationID *int64 `json:"github_installation_id,omitempty"`
	// Git repository owner
	GitRepositoryOwner *string `json:"git_repository_owner,omitempty"`
	// Git repository name
	GitRepository *string `json:"git_repository,omitempty"`
	// Kubernetes secret for this service
	KubernetesSecret string `json:"kubernetes_secret,omitempty"`
	// Reference the current active deployment
	CurrentDeploymentID *uuid.UUID `json:"current_deployment_id,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ServiceQuery when eager-loading is set.
	Edges        ServiceEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ServiceEdges holds the relations/edges for other nodes in the graph.
type ServiceEdges struct {
	// Environment holds the value of the environment edge.
	Environment *Environment `json:"environment,omitempty"`
	// GithubInstallation holds the value of the github_installation edge.
	GithubInstallation *GithubInstallation `json:"github_installation,omitempty"`
	// ServiceConfig holds the value of the service_config edge.
	ServiceConfig *ServiceConfig `json:"service_config,omitempty"`
	// Deployments holds the value of the deployments edge.
	Deployments []*Deployment `json:"deployments,omitempty"`
	// Optional reference to the currently active deployment
	CurrentDeployment *Deployment `json:"current_deployment,omitempty"`
	// VariableReferences holds the value of the variable_references edge.
	VariableReferences []*VariableReference `json:"variable_references,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [6]bool
}

// EnvironmentOrErr returns the Environment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEdges) EnvironmentOrErr() (*Environment, error) {
	if e.Environment != nil {
		return e.Environment, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: environment.Label}
	}
	return nil, &NotLoadedError{edge: "environment"}
}

// GithubInstallationOrErr returns the GithubInstallation value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEdges) GithubInstallationOrErr() (*GithubInstallation, error) {
	if e.GithubInstallation != nil {
		return e.GithubInstallation, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: githubinstallation.Label}
	}
	return nil, &NotLoadedError{edge: "github_installation"}
}

// ServiceConfigOrErr returns the ServiceConfig value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEdges) ServiceConfigOrErr() (*ServiceConfig, error) {
	if e.ServiceConfig != nil {
		return e.ServiceConfig, nil
	} else if e.loadedTypes[2] {
		return nil, &NotFoundError{label: serviceconfig.Label}
	}
	return nil, &NotLoadedError{edge: "service_config"}
}

// DeploymentsOrErr returns the Deployments value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) DeploymentsOrErr() ([]*Deployment, error) {
	if e.loadedTypes[3] {
		return e.Deployments, nil
	}
	return nil, &NotLoadedError{edge: "deployments"}
}

// CurrentDeploymentOrErr returns the CurrentDeployment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ServiceEdges) CurrentDeploymentOrErr() (*Deployment, error) {
	if e.CurrentDeployment != nil {
		return e.CurrentDeployment, nil
	} else if e.loadedTypes[4] {
		return nil, &NotFoundError{label: deployment.Label}
	}
	return nil, &NotLoadedError{edge: "current_deployment"}
}

// VariableReferencesOrErr returns the VariableReferences value or an error if the edge
// was not loaded in eager-loading.
func (e ServiceEdges) VariableReferencesOrErr() ([]*VariableReference, error) {
	if e.loadedTypes[5] {
		return e.VariableReferences, nil
	}
	return nil, &NotLoadedError{edge: "variable_references"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Service) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case service.FieldCurrentDeploymentID:
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		case service.FieldGithubInstallationID:
			values[i] = new(sql.NullInt64)
		case service.FieldName, service.FieldDisplayName, service.FieldDescription, service.FieldGitRepositoryOwner, service.FieldGitRepository, service.FieldKubernetesSecret:
			values[i] = new(sql.NullString)
		case service.FieldCreatedAt, service.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case service.FieldID, service.FieldEnvironmentID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Service fields.
func (s *Service) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case service.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				s.ID = *value
			}
		case service.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				s.CreatedAt = value.Time
			}
		case service.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				s.UpdatedAt = value.Time
			}
		case service.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				s.Name = value.String
			}
		case service.FieldDisplayName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field display_name", values[i])
			} else if value.Valid {
				s.DisplayName = value.String
			}
		case service.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				s.Description = value.String
			}
		case service.FieldEnvironmentID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field environment_id", values[i])
			} else if value != nil {
				s.EnvironmentID = *value
			}
		case service.FieldGithubInstallationID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field github_installation_id", values[i])
			} else if value.Valid {
				s.GithubInstallationID = new(int64)
				*s.GithubInstallationID = value.Int64
			}
		case service.FieldGitRepositoryOwner:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field git_repository_owner", values[i])
			} else if value.Valid {
				s.GitRepositoryOwner = new(string)
				*s.GitRepositoryOwner = value.String
			}
		case service.FieldGitRepository:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field git_repository", values[i])
			} else if value.Valid {
				s.GitRepository = new(string)
				*s.GitRepository = value.String
			}
		case service.FieldKubernetesSecret:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field kubernetes_secret", values[i])
			} else if value.Valid {
				s.KubernetesSecret = value.String
			}
		case service.FieldCurrentDeploymentID:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field current_deployment_id", values[i])
			} else if value.Valid {
				s.CurrentDeploymentID = new(uuid.UUID)
				*s.CurrentDeploymentID = *value.S.(*uuid.UUID)
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Service.
// This includes values selected through modifiers, order, etc.
func (s *Service) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryEnvironment queries the "environment" edge of the Service entity.
func (s *Service) QueryEnvironment() *EnvironmentQuery {
	return NewServiceClient(s.config).QueryEnvironment(s)
}

// QueryGithubInstallation queries the "github_installation" edge of the Service entity.
func (s *Service) QueryGithubInstallation() *GithubInstallationQuery {
	return NewServiceClient(s.config).QueryGithubInstallation(s)
}

// QueryServiceConfig queries the "service_config" edge of the Service entity.
func (s *Service) QueryServiceConfig() *ServiceConfigQuery {
	return NewServiceClient(s.config).QueryServiceConfig(s)
}

// QueryDeployments queries the "deployments" edge of the Service entity.
func (s *Service) QueryDeployments() *DeploymentQuery {
	return NewServiceClient(s.config).QueryDeployments(s)
}

// QueryCurrentDeployment queries the "current_deployment" edge of the Service entity.
func (s *Service) QueryCurrentDeployment() *DeploymentQuery {
	return NewServiceClient(s.config).QueryCurrentDeployment(s)
}

// QueryVariableReferences queries the "variable_references" edge of the Service entity.
func (s *Service) QueryVariableReferences() *VariableReferenceQuery {
	return NewServiceClient(s.config).QueryVariableReferences(s)
}

// Update returns a builder for updating this Service.
// Note that you need to call Service.Unwrap() before calling this method if this Service
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Service) Update() *ServiceUpdateOne {
	return NewServiceClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Service entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Service) Unwrap() *Service {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Service is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Service) String() string {
	var builder strings.Builder
	builder.WriteString("Service(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("created_at=")
	builder.WriteString(s.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(s.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(s.Name)
	builder.WriteString(", ")
	builder.WriteString("display_name=")
	builder.WriteString(s.DisplayName)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(s.Description)
	builder.WriteString(", ")
	builder.WriteString("environment_id=")
	builder.WriteString(fmt.Sprintf("%v", s.EnvironmentID))
	builder.WriteString(", ")
	if v := s.GithubInstallationID; v != nil {
		builder.WriteString("github_installation_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteString(", ")
	if v := s.GitRepositoryOwner; v != nil {
		builder.WriteString("git_repository_owner=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	if v := s.GitRepository; v != nil {
		builder.WriteString("git_repository=")
		builder.WriteString(*v)
	}
	builder.WriteString(", ")
	builder.WriteString("kubernetes_secret=")
	builder.WriteString(s.KubernetesSecret)
	builder.WriteString(", ")
	if v := s.CurrentDeploymentID; v != nil {
		builder.WriteString("current_deployment_id=")
		builder.WriteString(fmt.Sprintf("%v", *v))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Services is a parsable slice of Service.
type Services []*Service
