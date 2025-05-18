-- reverse: modify "user_teams" table
ALTER TABLE "user_teams" DROP CONSTRAINT "user_teams_user_id", DROP CONSTRAINT "user_teams_team_id";
-- reverse: modify "user_groups" table
ALTER TABLE "user_groups" DROP CONSTRAINT "user_groups_user_id", DROP CONSTRAINT "user_groups_group_id";
-- reverse: modify "group_permissions" table
ALTER TABLE "group_permissions" DROP CONSTRAINT "group_permissions_permission_id", DROP CONSTRAINT "group_permissions_group_id";
-- reverse: modify "webhooks" table
ALTER TABLE "webhooks" DROP CONSTRAINT "webhooks_teams_team_webhooks", DROP CONSTRAINT "webhooks_projects_project_webhooks";
-- reverse: modify "variable_references" table
ALTER TABLE "variable_references" DROP CONSTRAINT "variable_references_services_variable_references";
-- reverse: modify "service_groups" table
ALTER TABLE "service_groups" DROP CONSTRAINT "service_groups_environments_service_groups";
-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP CONSTRAINT "service_configs_services_service_config", DROP CONSTRAINT "service_configs_s3_sources_service_backup_source";
-- reverse: modify "services" table
ALTER TABLE "services" DROP CONSTRAINT "services_templates_services", DROP CONSTRAINT "services_service_groups_services", DROP CONSTRAINT "services_github_installations_services", DROP CONSTRAINT "services_environments_services", DROP CONSTRAINT "services_deployments_current_deployment";
-- reverse: modify "s3_sources" table
ALTER TABLE "s3_sources" DROP CONSTRAINT "s3_sources_teams_s3_sources";
-- reverse: modify "projects" table
ALTER TABLE "projects" DROP CONSTRAINT "projects_teams_projects", DROP CONSTRAINT "projects_environments_project_default";
-- reverse: modify "oauth2_tokens" table
ALTER TABLE "oauth2_tokens" DROP CONSTRAINT "oauth2_tokens_users_oauth2_tokens";
-- reverse: modify "oauth2_codes" table
ALTER TABLE "oauth2_codes" DROP CONSTRAINT "oauth2_codes_users_oauth2_codes";
-- reverse: modify "github_installations" table
ALTER TABLE "github_installations" DROP CONSTRAINT "github_installations_github_apps_installations";
-- reverse: modify "github_apps" table
ALTER TABLE "github_apps" DROP CONSTRAINT "github_apps_users_created_by";
-- reverse: modify "environments" table
ALTER TABLE "environments" DROP CONSTRAINT "environments_projects_environments";
-- reverse: modify "deployments" table
ALTER TABLE "deployments" DROP CONSTRAINT "deployments_services_deployments";
-- reverse: create "user_teams" table
DROP TABLE "user_teams";
-- reverse: create "user_groups" table
DROP TABLE "user_groups";
-- reverse: create "group_permissions" table
DROP TABLE "group_permissions";
-- reverse: create "webhooks" table
DROP TABLE "webhooks";
-- reverse: create index "variablereference_target_service_id_target_name" to table: "variable_references"
DROP INDEX "variablereference_target_service_id_target_name";
-- reverse: create "variable_references" table
DROP TABLE "variable_references";
-- reverse: create index "users_email_key" to table: "users"
DROP INDEX "users_email_key";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create index "template_name_version" to table: "templates"
DROP INDEX "template_name_version";
-- reverse: create "templates" table
DROP TABLE "templates";
-- reverse: create index "teams_namespace_key" to table: "teams"
DROP INDEX "teams_namespace_key";
-- reverse: create index "teams_kubernetes_name_key" to table: "teams"
DROP INDEX "teams_kubernetes_name_key";
-- reverse: create "teams" table
DROP TABLE "teams";
-- reverse: create "system_settings" table
DROP TABLE "system_settings";
-- reverse: create "service_groups" table
DROP TABLE "service_groups";
-- reverse: create index "service_configs_service_id_key" to table: "service_configs"
DROP INDEX "service_configs_service_id_key";
-- reverse: create "service_configs" table
DROP TABLE "service_configs";
-- reverse: create index "services_kubernetes_name_key" to table: "services"
DROP INDEX "services_kubernetes_name_key";
-- reverse: create "services" table
DROP TABLE "services";
-- reverse: create "s3_sources" table
DROP TABLE "s3_sources";
-- reverse: create "registries" table
DROP TABLE "registries";
-- reverse: create index "projects_kubernetes_name_key" to table: "projects"
DROP INDEX "projects_kubernetes_name_key";
-- reverse: create "projects" table
DROP TABLE "projects";
-- reverse: create "permissions" table
DROP TABLE "permissions";
-- reverse: create index "oauth2_tokens_refresh_token_key" to table: "oauth2_tokens"
DROP INDEX "oauth2_tokens_refresh_token_key";
-- reverse: create index "oauth2_tokens_access_token_key" to table: "oauth2_tokens"
DROP INDEX "oauth2_tokens_access_token_key";
-- reverse: create "oauth2_tokens" table
DROP TABLE "oauth2_tokens";
-- reverse: create index "oauth2_codes_auth_code_key" to table: "oauth2_codes"
DROP INDEX "oauth2_codes_auth_code_key";
-- reverse: create "oauth2_codes" table
DROP TABLE "oauth2_codes";
-- reverse: create "jwt_keys" table
DROP TABLE "jwt_keys";
-- reverse: create "groups" table
DROP TABLE "groups";
-- reverse: create index "githubinstallation_github_app_id" to table: "github_installations"
DROP INDEX "githubinstallation_github_app_id";
-- reverse: create "github_installations" table
DROP TABLE "github_installations";
-- reverse: create index "github_apps_uuid_key" to table: "github_apps"
DROP INDEX "github_apps_uuid_key";
-- reverse: create "github_apps" table
DROP TABLE "github_apps";
-- reverse: create index "environments_kubernetes_name_key" to table: "environments"
DROP INDEX "environments_kubernetes_name_key";
-- reverse: create "environments" table
DROP TABLE "environments";
-- reverse: create "deployments" table
DROP TABLE "deployments";
-- reverse: create "bootstrap_flag" table
DROP TABLE "bootstrap_flag";
