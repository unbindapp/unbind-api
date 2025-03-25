// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// DeploymentsColumns holds the columns for the "deployments" table.
	DeploymentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "status", Type: field.TypeEnum, Enums: []string{"queued", "building", "succeeded", "cancelled", "failed"}},
		{Name: "source", Type: field.TypeEnum, Enums: []string{"manual", "git"}, Default: "manual"},
		{Name: "error", Type: field.TypeString, Nullable: true},
		{Name: "commit_sha", Type: field.TypeString, Nullable: true},
		{Name: "commit_message", Type: field.TypeString, Nullable: true},
		{Name: "started_at", Type: field.TypeTime, Nullable: true},
		{Name: "completed_at", Type: field.TypeTime, Nullable: true},
		{Name: "kubernetes_job_name", Type: field.TypeString, Nullable: true},
		{Name: "kubernetes_job_status", Type: field.TypeString, Nullable: true},
		{Name: "attempts", Type: field.TypeInt, Default: 0},
		{Name: "service_id", Type: field.TypeUUID},
	}
	// DeploymentsTable holds the schema information for the "deployments" table.
	DeploymentsTable = &schema.Table{
		Name:       "deployments",
		Columns:    DeploymentsColumns,
		PrimaryKey: []*schema.Column{DeploymentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "deployments_services_deployments",
				Columns:    []*schema.Column{DeploymentsColumns[13]},
				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// EnvironmentsColumns holds the columns for the "environments" table.
	EnvironmentsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "display_name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "kubernetes_secret", Type: field.TypeString},
		{Name: "project_id", Type: field.TypeUUID},
	}
	// EnvironmentsTable holds the schema information for the "environments" table.
	EnvironmentsTable = &schema.Table{
		Name:       "environments",
		Columns:    EnvironmentsColumns,
		PrimaryKey: []*schema.Column{EnvironmentsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "environments_projects_environments",
				Columns:    []*schema.Column{EnvironmentsColumns[8]},
				RefColumns: []*schema.Column{ProjectsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// GithubAppsColumns holds the columns for the "github_apps" table.
	GithubAppsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "client_id", Type: field.TypeString},
		{Name: "client_secret", Type: field.TypeString},
		{Name: "webhook_secret", Type: field.TypeString},
		{Name: "private_key", Type: field.TypeString, Size: 2147483647},
		{Name: "created_by", Type: field.TypeUUID},
	}
	// GithubAppsTable holds the schema information for the "github_apps" table.
	GithubAppsTable = &schema.Table{
		Name:       "github_apps",
		Columns:    GithubAppsColumns,
		PrimaryKey: []*schema.Column{GithubAppsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "github_apps_users_created_by",
				Columns:    []*schema.Column{GithubAppsColumns[8]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// GithubInstallationsColumns holds the columns for the "github_installations" table.
	GithubInstallationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "account_id", Type: field.TypeInt64},
		{Name: "account_login", Type: field.TypeString},
		{Name: "account_type", Type: field.TypeEnum, Enums: []string{"Organization", "User"}},
		{Name: "account_url", Type: field.TypeString},
		{Name: "repository_selection", Type: field.TypeEnum, Enums: []string{"all", "selected"}, Default: "all"},
		{Name: "suspended", Type: field.TypeBool, Default: false},
		{Name: "active", Type: field.TypeBool, Default: true},
		{Name: "permissions", Type: field.TypeJSON, Nullable: true},
		{Name: "events", Type: field.TypeJSON, Nullable: true},
		{Name: "github_app_id", Type: field.TypeInt64},
	}
	// GithubInstallationsTable holds the schema information for the "github_installations" table.
	GithubInstallationsTable = &schema.Table{
		Name:       "github_installations",
		Columns:    GithubInstallationsColumns,
		PrimaryKey: []*schema.Column{GithubInstallationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "github_installations_github_apps_installations",
				Columns:    []*schema.Column{GithubInstallationsColumns[12]},
				RefColumns: []*schema.Column{GithubAppsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "githubinstallation_github_app_id",
				Unique:  true,
				Columns: []*schema.Column{GithubInstallationsColumns[12]},
			},
		},
	}
	// GroupsColumns holds the columns for the "groups" table.
	GroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "superuser", Type: field.TypeBool, Default: false},
		{Name: "k8s_role_name", Type: field.TypeString, Nullable: true},
		{Name: "identity_provider", Type: field.TypeString, Nullable: true},
		{Name: "external_id", Type: field.TypeString, Nullable: true},
		{Name: "team_id", Type: field.TypeUUID, Nullable: true},
	}
	// GroupsTable holds the schema information for the "groups" table.
	GroupsTable = &schema.Table{
		Name:       "groups",
		Columns:    GroupsColumns,
		PrimaryKey: []*schema.Column{GroupsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "groups_teams_groups",
				Columns:    []*schema.Column{GroupsColumns[9]},
				RefColumns: []*schema.Column{TeamsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "group_name_team_id",
				Unique:  true,
				Columns: []*schema.Column{GroupsColumns[3], GroupsColumns[9]},
			},
		},
	}
	// JwtKeysColumns holds the columns for the "jwt_keys" table.
	JwtKeysColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "label", Type: field.TypeString},
		{Name: "private_key", Type: field.TypeBytes},
	}
	// JwtKeysTable holds the schema information for the "jwt_keys" table.
	JwtKeysTable = &schema.Table{
		Name:       "jwt_keys",
		Columns:    JwtKeysColumns,
		PrimaryKey: []*schema.Column{JwtKeysColumns[0]},
	}
	// Oauth2CodesColumns holds the columns for the "oauth2_codes" table.
	Oauth2CodesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "auth_code", Type: field.TypeString, Unique: true},
		{Name: "client_id", Type: field.TypeString},
		{Name: "scope", Type: field.TypeString},
		{Name: "expires_at", Type: field.TypeTime},
		{Name: "revoked", Type: field.TypeBool, Default: false},
		{Name: "user_oauth2_codes", Type: field.TypeUUID},
	}
	// Oauth2CodesTable holds the schema information for the "oauth2_codes" table.
	Oauth2CodesTable = &schema.Table{
		Name:       "oauth2_codes",
		Columns:    Oauth2CodesColumns,
		PrimaryKey: []*schema.Column{Oauth2CodesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "oauth2_codes_users_oauth2_codes",
				Columns:    []*schema.Column{Oauth2CodesColumns[8]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// Oauth2TokensColumns holds the columns for the "oauth2_tokens" table.
	Oauth2TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "access_token", Type: field.TypeString, Unique: true},
		{Name: "refresh_token", Type: field.TypeString, Unique: true},
		{Name: "client_id", Type: field.TypeString},
		{Name: "expires_at", Type: field.TypeTime},
		{Name: "revoked", Type: field.TypeBool, Default: false},
		{Name: "scope", Type: field.TypeString},
		{Name: "device_info", Type: field.TypeString, Nullable: true},
		{Name: "user_oauth2_tokens", Type: field.TypeUUID},
	}
	// Oauth2TokensTable holds the schema information for the "oauth2_tokens" table.
	Oauth2TokensTable = &schema.Table{
		Name:       "oauth2_tokens",
		Columns:    Oauth2TokensColumns,
		PrimaryKey: []*schema.Column{Oauth2TokensColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "oauth2_tokens_users_oauth2_tokens",
				Columns:    []*schema.Column{Oauth2TokensColumns[10]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// PermissionsColumns holds the columns for the "permissions" table.
	PermissionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "action", Type: field.TypeEnum, Enums: []string{"read", "create", "update", "delete", "manage", "admin", "edit", "view"}},
		{Name: "resource_type", Type: field.TypeEnum, Enums: []string{"team", "project", "group", "environment", "permission", "user", "system", "service"}},
		{Name: "resource_id", Type: field.TypeString},
		{Name: "scope", Type: field.TypeString, Nullable: true},
		{Name: "labels", Type: field.TypeJSON, Nullable: true},
	}
	// PermissionsTable holds the schema information for the "permissions" table.
	PermissionsTable = &schema.Table{
		Name:       "permissions",
		Columns:    PermissionsColumns,
		PrimaryKey: []*schema.Column{PermissionsColumns[0]},
	}
	// ProjectsColumns holds the columns for the "projects" table.
	ProjectsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "display_name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "status", Type: field.TypeString, Default: "active"},
		{Name: "kubernetes_secret", Type: field.TypeString},
		{Name: "team_id", Type: field.TypeUUID},
	}
	// ProjectsTable holds the schema information for the "projects" table.
	ProjectsTable = &schema.Table{
		Name:       "projects",
		Columns:    ProjectsColumns,
		PrimaryKey: []*schema.Column{ProjectsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "projects_teams_projects",
				Columns:    []*schema.Column{ProjectsColumns[8]},
				RefColumns: []*schema.Column{TeamsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// ServicesColumns holds the columns for the "services" table.
	ServicesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString},
		{Name: "display_name", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
		{Name: "git_repository_owner", Type: field.TypeString, Nullable: true},
		{Name: "git_repository", Type: field.TypeString, Nullable: true},
		{Name: "kubernetes_secret", Type: field.TypeString},
		{Name: "environment_id", Type: field.TypeUUID},
		{Name: "github_installation_id", Type: field.TypeInt64, Nullable: true},
	}
	// ServicesTable holds the schema information for the "services" table.
	ServicesTable = &schema.Table{
		Name:       "services",
		Columns:    ServicesColumns,
		PrimaryKey: []*schema.Column{ServicesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "services_environments_services",
				Columns:    []*schema.Column{ServicesColumns[9]},
				RefColumns: []*schema.Column{EnvironmentsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "services_github_installations_services",
				Columns:    []*schema.Column{ServicesColumns[10]},
				RefColumns: []*schema.Column{GithubInstallationsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// ServiceConfigsColumns holds the columns for the "service_configs" table.
	ServiceConfigsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "type", Type: field.TypeEnum, Enums: []string{"github", "docker-image"}},
		{Name: "builder", Type: field.TypeEnum, Enums: []string{"nixpacks", "railpack", "docker"}},
		{Name: "provider", Type: field.TypeEnum, Nullable: true, Enums: []string{"node", "deno", "go", "java", "php", "python", "staticfile", "unknown"}},
		{Name: "framework", Type: field.TypeEnum, Nullable: true, Enums: []string{"next", "astro", "vite", "cra", "angular", "remix", "bun", "express", "python", "django", "flask", "fastapi", "fasthtml", "gin", "spring-boot", "laravel", "unknown"}},
		{Name: "git_branch", Type: field.TypeString, Nullable: true},
		{Name: "hosts", Type: field.TypeJSON, Nullable: true},
		{Name: "ports", Type: field.TypeJSON, Nullable: true},
		{Name: "replicas", Type: field.TypeInt32, Default: 2},
		{Name: "auto_deploy", Type: field.TypeBool, Default: false},
		{Name: "run_command", Type: field.TypeString, Nullable: true},
		{Name: "public", Type: field.TypeBool, Default: false},
		{Name: "image", Type: field.TypeString, Nullable: true},
		{Name: "service_id", Type: field.TypeUUID, Unique: true},
	}
	// ServiceConfigsTable holds the schema information for the "service_configs" table.
	ServiceConfigsTable = &schema.Table{
		Name:       "service_configs",
		Columns:    ServiceConfigsColumns,
		PrimaryKey: []*schema.Column{ServiceConfigsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "service_configs_services_service_config",
				Columns:    []*schema.Column{ServiceConfigsColumns[15]},
				RefColumns: []*schema.Column{ServicesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// TeamsColumns holds the columns for the "teams" table.
	TeamsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "display_name", Type: field.TypeString},
		{Name: "namespace", Type: field.TypeString, Unique: true},
		{Name: "kubernetes_secret", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Nullable: true},
	}
	// TeamsTable holds the schema information for the "teams" table.
	TeamsTable = &schema.Table{
		Name:       "teams",
		Columns:    TeamsColumns,
		PrimaryKey: []*schema.Column{TeamsColumns[0]},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "password_hash", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// GroupPermissionsColumns holds the columns for the "group_permissions" table.
	GroupPermissionsColumns = []*schema.Column{
		{Name: "group_id", Type: field.TypeUUID},
		{Name: "permission_id", Type: field.TypeUUID},
	}
	// GroupPermissionsTable holds the schema information for the "group_permissions" table.
	GroupPermissionsTable = &schema.Table{
		Name:       "group_permissions",
		Columns:    GroupPermissionsColumns,
		PrimaryKey: []*schema.Column{GroupPermissionsColumns[0], GroupPermissionsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "group_permissions_group_id",
				Columns:    []*schema.Column{GroupPermissionsColumns[0]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "group_permissions_permission_id",
				Columns:    []*schema.Column{GroupPermissionsColumns[1]},
				RefColumns: []*schema.Column{PermissionsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UserGroupsColumns holds the columns for the "user_groups" table.
	UserGroupsColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "group_id", Type: field.TypeUUID},
	}
	// UserGroupsTable holds the schema information for the "user_groups" table.
	UserGroupsTable = &schema.Table{
		Name:       "user_groups",
		Columns:    UserGroupsColumns,
		PrimaryKey: []*schema.Column{UserGroupsColumns[0], UserGroupsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_groups_user_id",
				Columns:    []*schema.Column{UserGroupsColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_groups_group_id",
				Columns:    []*schema.Column{UserGroupsColumns[1]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UserTeamsColumns holds the columns for the "user_teams" table.
	UserTeamsColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeUUID},
		{Name: "team_id", Type: field.TypeUUID},
	}
	// UserTeamsTable holds the schema information for the "user_teams" table.
	UserTeamsTable = &schema.Table{
		Name:       "user_teams",
		Columns:    UserTeamsColumns,
		PrimaryKey: []*schema.Column{UserTeamsColumns[0], UserTeamsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_teams_user_id",
				Columns:    []*schema.Column{UserTeamsColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_teams_team_id",
				Columns:    []*schema.Column{UserTeamsColumns[1]},
				RefColumns: []*schema.Column{TeamsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		DeploymentsTable,
		EnvironmentsTable,
		GithubAppsTable,
		GithubInstallationsTable,
		GroupsTable,
		JwtKeysTable,
		Oauth2CodesTable,
		Oauth2TokensTable,
		PermissionsTable,
		ProjectsTable,
		ServicesTable,
		ServiceConfigsTable,
		TeamsTable,
		UsersTable,
		GroupPermissionsTable,
		UserGroupsTable,
		UserTeamsTable,
	}
)

func init() {
	DeploymentsTable.ForeignKeys[0].RefTable = ServicesTable
	DeploymentsTable.Annotation = &entsql.Annotation{
		Table: "deployments",
	}
	EnvironmentsTable.ForeignKeys[0].RefTable = ProjectsTable
	EnvironmentsTable.Annotation = &entsql.Annotation{
		Table: "environments",
	}
	GithubAppsTable.ForeignKeys[0].RefTable = UsersTable
	GithubAppsTable.Annotation = &entsql.Annotation{
		Table: "github_apps",
	}
	GithubInstallationsTable.ForeignKeys[0].RefTable = GithubAppsTable
	GithubInstallationsTable.Annotation = &entsql.Annotation{
		Table: "github_installations",
	}
	GroupsTable.ForeignKeys[0].RefTable = TeamsTable
	GroupsTable.Annotation = &entsql.Annotation{
		Table: "groups",
	}
	JwtKeysTable.Annotation = &entsql.Annotation{
		Table: "jwt_keys",
	}
	Oauth2CodesTable.ForeignKeys[0].RefTable = UsersTable
	Oauth2CodesTable.Annotation = &entsql.Annotation{
		Table: "oauth2_codes",
	}
	Oauth2TokensTable.ForeignKeys[0].RefTable = UsersTable
	Oauth2TokensTable.Annotation = &entsql.Annotation{
		Table: "oauth2_tokens",
	}
	PermissionsTable.Annotation = &entsql.Annotation{
		Table: "permissions",
	}
	ProjectsTable.ForeignKeys[0].RefTable = TeamsTable
	ProjectsTable.Annotation = &entsql.Annotation{
		Table: "projects",
	}
	ServicesTable.ForeignKeys[0].RefTable = EnvironmentsTable
	ServicesTable.ForeignKeys[1].RefTable = GithubInstallationsTable
	ServicesTable.Annotation = &entsql.Annotation{
		Table: "services",
	}
	ServiceConfigsTable.ForeignKeys[0].RefTable = ServicesTable
	ServiceConfigsTable.Annotation = &entsql.Annotation{
		Table: "service_configs",
	}
	TeamsTable.Annotation = &entsql.Annotation{
		Table: "teams",
	}
	UsersTable.Annotation = &entsql.Annotation{
		Table: "users",
	}
	GroupPermissionsTable.ForeignKeys[0].RefTable = GroupsTable
	GroupPermissionsTable.ForeignKeys[1].RefTable = PermissionsTable
	UserGroupsTable.ForeignKeys[0].RefTable = UsersTable
	UserGroupsTable.ForeignKeys[1].RefTable = GroupsTable
	UserTeamsTable.ForeignKeys[0].RefTable = UsersTable
	UserTeamsTable.ForeignKeys[1].RefTable = TeamsTable
}
