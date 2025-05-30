-- +goose Up
-- create "bootstrap_flag" table
CREATE TABLE "bootstrap_flag" (
  "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
  "is_bootstrapped" boolean NOT NULL,
  PRIMARY KEY ("id")
);
-- create "deployments" table
CREATE TABLE "deployments" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "status" character varying NOT NULL,
  "source" character varying NOT NULL DEFAULT 'manual',
  "error" character varying NULL,
  "commit_sha" character varying NULL,
  "commit_message" character varying NULL,
  "commit_author" jsonb NULL,
  "queued_at" timestamptz NULL,
  "started_at" timestamptz NULL,
  "completed_at" timestamptz NULL,
  "kubernetes_job_name" character varying NULL,
  "kubernetes_job_status" character varying NULL,
  "attempts" bigint NOT NULL DEFAULT 0,
  "image" character varying NULL,
  "resource_definition" jsonb NULL,
  "service_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create "environments" table
CREATE TABLE "environments" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "kubernetes_name" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "active" boolean NOT NULL DEFAULT true,
  "kubernetes_secret" character varying NOT NULL,
  "project_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "environments_kubernetes_name_key" to table: "environments"
CREATE UNIQUE INDEX "environments_kubernetes_name_key" ON "environments" ("kubernetes_name");
-- create "github_apps" table
CREATE TABLE "github_apps" (
  "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "uuid" uuid NOT NULL,
  "name" character varying NOT NULL,
  "client_id" character varying NOT NULL,
  "client_secret" character varying NOT NULL,
  "webhook_secret" character varying NOT NULL,
  "private_key" text NOT NULL,
  "created_by" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "github_apps_uuid_key" to table: "github_apps"
CREATE UNIQUE INDEX "github_apps_uuid_key" ON "github_apps" ("uuid");
-- create "github_installations" table
CREATE TABLE "github_installations" (
  "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "account_id" bigint NOT NULL,
  "account_login" character varying NOT NULL,
  "account_type" character varying NOT NULL,
  "account_url" character varying NOT NULL,
  "repository_selection" character varying NOT NULL DEFAULT 'all',
  "suspended" boolean NOT NULL DEFAULT false,
  "active" boolean NOT NULL DEFAULT true,
  "permissions" jsonb NULL,
  "events" jsonb NULL,
  "github_app_id" bigint NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "githubinstallation_github_app_id" to table: "github_installations"
CREATE UNIQUE INDEX "githubinstallation_github_app_id" ON "github_installations" ("github_app_id");
-- create "groups" table
CREATE TABLE "groups" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "k8s_role_name" character varying NULL,
  PRIMARY KEY ("id")
);
-- create "jwt_keys" table
CREATE TABLE "jwt_keys" (
  "id" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
  "label" character varying NOT NULL,
  "private_key" bytea NOT NULL,
  PRIMARY KEY ("id")
);
-- create "oauth2_codes" table
CREATE TABLE "oauth2_codes" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "auth_code" character varying NOT NULL,
  "client_id" character varying NOT NULL,
  "scope" character varying NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked" boolean NOT NULL DEFAULT false,
  "user_oauth2_codes" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "oauth2_codes_auth_code_key" to table: "oauth2_codes"
CREATE UNIQUE INDEX "oauth2_codes_auth_code_key" ON "oauth2_codes" ("auth_code");
-- create "oauth2_tokens" table
CREATE TABLE "oauth2_tokens" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "access_token" character varying NOT NULL,
  "refresh_token" character varying NOT NULL,
  "client_id" character varying NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked" boolean NOT NULL DEFAULT false,
  "scope" character varying NOT NULL,
  "device_info" character varying NULL,
  "user_oauth2_tokens" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "oauth2_tokens_access_token_key" to table: "oauth2_tokens"
CREATE UNIQUE INDEX "oauth2_tokens_access_token_key" ON "oauth2_tokens" ("access_token");
-- create index "oauth2_tokens_refresh_token_key" to table: "oauth2_tokens"
CREATE UNIQUE INDEX "oauth2_tokens_refresh_token_key" ON "oauth2_tokens" ("refresh_token");
-- create "permissions" table
CREATE TABLE "permissions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "action" character varying NOT NULL,
  "resource_type" character varying NOT NULL,
  "resource_selector" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create "projects" table
CREATE TABLE "projects" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "kubernetes_name" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "status" character varying NOT NULL DEFAULT 'active',
  "kubernetes_secret" character varying NOT NULL,
  "default_environment_id" uuid NULL,
  "team_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "projects_kubernetes_name_key" to table: "projects"
CREATE UNIQUE INDEX "projects_kubernetes_name_key" ON "projects" ("kubernetes_name");
-- create "registries" table
CREATE TABLE "registries" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "host" character varying NOT NULL,
  "kubernetes_secret" character varying NULL,
  "is_default" boolean NOT NULL,
  PRIMARY KEY ("id")
);
-- create "s3_sources" table
CREATE TABLE "s3_sources" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "endpoint" character varying NOT NULL,
  "region" character varying NOT NULL,
  "force_path_style" boolean NOT NULL DEFAULT true,
  "kubernetes_secret" character varying NOT NULL,
  "team_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create "services" table
CREATE TABLE "services" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "type" character varying NOT NULL,
  "kubernetes_name" character varying NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NULL,
  "database" character varying NULL,
  "database_version" character varying NULL,
  "git_repository_owner" character varying NULL,
  "git_repository" character varying NULL,
  "kubernetes_secret" character varying NOT NULL,
  "template_instance_id" uuid NULL,
  "environment_id" uuid NOT NULL,
  "github_installation_id" bigint NULL,
  "current_deployment_id" uuid NULL,
  "service_group_id" uuid NULL,
  "template_id" uuid NULL,
  PRIMARY KEY ("id")
);
-- create index "services_kubernetes_name_key" to table: "services"
CREATE UNIQUE INDEX "services_kubernetes_name_key" ON "services" ("kubernetes_name");
-- create "service_configs" table
CREATE TABLE "service_configs" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "builder" character varying NOT NULL,
  "icon" character varying NOT NULL,
  "dockerfile_path" character varying NULL,
  "dockerfile_context" character varying NULL,
  "railpack_provider" character varying NULL,
  "railpack_framework" character varying NULL,
  "git_branch" character varying NULL,
  "git_tag" character varying NULL,
  "hosts" jsonb NULL,
  "ports" jsonb NULL,
  "replicas" integer NOT NULL DEFAULT 1,
  "auto_deploy" boolean NOT NULL DEFAULT false,
  "install_command" character varying NULL,
  "build_command" character varying NULL,
  "run_command" character varying NULL,
  "is_public" boolean NOT NULL DEFAULT false,
  "image" character varying NULL,
  "definition_version" character varying NULL,
  "database_config" jsonb NULL,
  "s3_backup_bucket" character varying NULL,
  "backup_schedule" character varying NOT NULL DEFAULT '5 5 * * *',
  "backup_retention_count" bigint NOT NULL DEFAULT 3,
  "volumes" jsonb NULL,
  "security_context" jsonb NULL,
  "health_check" jsonb NULL,
  "variable_mounts" jsonb NULL,
  "protected_variables" jsonb NULL,
  "s3_backup_source_id" uuid NULL,
  "service_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "service_configs_service_id_key" to table: "service_configs"
CREATE UNIQUE INDEX "service_configs_service_id_key" ON "service_configs" ("service_id");
-- create "service_groups" table
CREATE TABLE "service_groups" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "icon" character varying NULL,
  "description" character varying NULL,
  "environment_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create "system_settings" table
CREATE TABLE "system_settings" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "wildcard_base_url" character varying NULL,
  "buildkit_settings" jsonb NULL,
  PRIMARY KEY ("id")
);
-- create "teams" table
CREATE TABLE "teams" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "kubernetes_name" character varying NOT NULL,
  "name" character varying NOT NULL,
  "namespace" character varying NOT NULL,
  "kubernetes_secret" character varying NOT NULL,
  "description" character varying NULL,
  PRIMARY KEY ("id")
);
-- create index "teams_kubernetes_name_key" to table: "teams"
CREATE UNIQUE INDEX "teams_kubernetes_name_key" ON "teams" ("kubernetes_name");
-- create index "teams_namespace_key" to table: "teams"
CREATE UNIQUE INDEX "teams_namespace_key" ON "teams" ("namespace");
-- create "templates" table
CREATE TABLE "templates" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "description" character varying NOT NULL,
  "icon" character varying NOT NULL,
  "keywords" jsonb NULL,
  "display_rank" bigint NOT NULL DEFAULT 0,
  "version" bigint NOT NULL,
  "immutable" boolean NOT NULL DEFAULT false,
  "definition" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "template_name_version" to table: "templates"
CREATE UNIQUE INDEX "template_name_version" ON "templates" ("name", "version");
-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "email" character varying NOT NULL,
  "password_hash" character varying NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "users" ("email");
-- create "variable_references" table
CREATE TABLE "variable_references" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "target_name" character varying NOT NULL,
  "sources" jsonb NOT NULL,
  "value_template" character varying NOT NULL,
  "error" character varying NULL,
  "target_service_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "variablereference_target_service_id_target_name" to table: "variable_references"
CREATE UNIQUE INDEX "variablereference_target_service_id_target_name" ON "variable_references" ("target_service_id", "target_name");
-- create "webhooks" table
CREATE TABLE "webhooks" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "url" character varying NOT NULL,
  "type" character varying NOT NULL,
  "events" jsonb NOT NULL,
  "project_id" uuid NULL,
  "team_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);
-- create "group_permissions" table
CREATE TABLE "group_permissions" (
  "group_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  PRIMARY KEY ("group_id", "permission_id")
);
-- create "user_groups" table
CREATE TABLE "user_groups" (
  "user_id" uuid NOT NULL,
  "group_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "group_id")
);
-- create "user_teams" table
CREATE TABLE "user_teams" (
  "user_id" uuid NOT NULL,
  "team_id" uuid NOT NULL,
  PRIMARY KEY ("user_id", "team_id")
);
-- modify "deployments" table
ALTER TABLE "deployments" ADD
CONSTRAINT "deployments_services_deployments" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "environments" table
ALTER TABLE "environments" ADD
CONSTRAINT "environments_projects_environments" FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "github_apps" table
ALTER TABLE "github_apps" ADD
CONSTRAINT "github_apps_users_created_by" FOREIGN KEY ("created_by") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "github_installations" table
ALTER TABLE "github_installations" ADD
CONSTRAINT "github_installations_github_apps_installations" FOREIGN KEY ("github_app_id") REFERENCES "github_apps" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "oauth2_codes" table
ALTER TABLE "oauth2_codes" ADD
CONSTRAINT "oauth2_codes_users_oauth2_codes" FOREIGN KEY ("user_oauth2_codes") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "oauth2_tokens" table
ALTER TABLE "oauth2_tokens" ADD
CONSTRAINT "oauth2_tokens_users_oauth2_tokens" FOREIGN KEY ("user_oauth2_tokens") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "projects" table
ALTER TABLE "projects" ADD
CONSTRAINT "projects_environments_project_default" FOREIGN KEY ("default_environment_id") REFERENCES "environments" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD
CONSTRAINT "projects_teams_projects" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "s3_sources" table
ALTER TABLE "s3_sources" ADD
CONSTRAINT "s3_sources_teams_s3_sources" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "services" table
ALTER TABLE "services" ADD
CONSTRAINT "services_deployments_current_deployment" FOREIGN KEY ("current_deployment_id") REFERENCES "deployments" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD
CONSTRAINT "services_environments_services" FOREIGN KEY ("environment_id") REFERENCES "environments" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
CONSTRAINT "services_github_installations_services" FOREIGN KEY ("github_installation_id") REFERENCES "github_installations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD
CONSTRAINT "services_service_groups_services" FOREIGN KEY ("service_group_id") REFERENCES "service_groups" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD
CONSTRAINT "services_templates_services" FOREIGN KEY ("template_id") REFERENCES "templates" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- modify "service_configs" table
ALTER TABLE "service_configs" ADD
CONSTRAINT "service_configs_s3_sources_service_backup_source" FOREIGN KEY ("s3_backup_source_id") REFERENCES "s3_sources" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD
CONSTRAINT "service_configs_services_service_config" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "service_groups" table
ALTER TABLE "service_groups" ADD
CONSTRAINT "service_groups_environments_service_groups" FOREIGN KEY ("environment_id") REFERENCES "environments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "variable_references" table
ALTER TABLE "variable_references" ADD
CONSTRAINT "variable_references_services_variable_references" FOREIGN KEY ("target_service_id") REFERENCES "services" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "webhooks" table
ALTER TABLE "webhooks" ADD
CONSTRAINT "webhooks_projects_project_webhooks" FOREIGN KEY ("project_id") REFERENCES "projects" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
CONSTRAINT "webhooks_teams_team_webhooks" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "group_permissions" table
ALTER TABLE "group_permissions" ADD
CONSTRAINT "group_permissions_group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
CONSTRAINT "group_permissions_permission_id" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "user_groups" table
ALTER TABLE "user_groups" ADD
CONSTRAINT "user_groups_group_id" FOREIGN KEY ("group_id") REFERENCES "groups" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
CONSTRAINT "user_groups_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- modify "user_teams" table
ALTER TABLE "user_teams" ADD
CONSTRAINT "user_teams_team_id" FOREIGN KEY ("team_id") REFERENCES "teams" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
CONSTRAINT "user_teams_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;

-- +goose Down
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
