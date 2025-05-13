package templates

import (
	"github.com/unbindapp/unbind-api/ent/schema"
	"github.com/unbindapp/unbind-api/internal/common/utils"
)

// supabaseTemplate returns a template definition for Supabase
func supabaseTemplate() *schema.TemplateDefinition {
	return &schema.TemplateDefinition{
		Name:        "supabase",
		Description: "The open source Firebase alternative",
		Version:     1,
		Inputs: []schema.TemplateInput{
			{
				ID:          1,
				Name:        "API Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Supabase API.",
				Required:    true,
				TargetPort:  utils.ToPtr(8000),
			},
			{
				ID:          2,
				Name:        "Studio Host",
				Type:        schema.InputTypeHost,
				Description: "Hostname to use for the Supabase Studio.",
				Required:    true,
				TargetPort:  utils.ToPtr(3000),
			},
			{
				ID:          3,
				Name:        "Storage Size",
				Type:        schema.InputTypeVolumeSize,
				Description: "Size of the persistent storage for Supabase storage service.",
				Required:    true,
				Default:     utils.ToPtr("10Gi"),
			},
		},
		Services: []schema.TemplateService{
			{
				ID:           1,
				Name:         "PostgreSQL",
				Type:         schema.ServiceTypeDatabase,
				Builder:      schema.ServiceBuilderDatabase,
				DatabaseType: utils.ToPtr("postgres"),
				DatabaseConfig: &schema.DatabaseConfig{
					DefaultDatabaseName: "postgres",
					Version:             "15",
					InitDB: `
-- Create _supabase database
CREATE DATABASE _supabase WITH OWNER postgres;

-- Create _realtime schema
\c _supabase
create schema if not exists _realtime;
alter schema _realtime owner to postgres;

-- Create _supavisor schema
create schema if not exists _supavisor;
alter schema _supavisor owner to postgres;

-- Create _analytics schema
create schema if not exists _analytics;
alter schema _analytics owner to postgres;

-- Create pg_net extension and supabase_functions schema
\c postgres
CREATE EXTENSION IF NOT EXISTS pg_net SCHEMA extensions;

-- Create supabase_functions schema
CREATE SCHEMA supabase_functions AUTHORIZATION postgres;
GRANT USAGE ON SCHEMA supabase_functions TO postgres, anon, authenticated, service_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA supabase_functions GRANT ALL ON TABLES TO postgres, anon, authenticated, service_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA supabase_functions GRANT ALL ON FUNCTIONS TO postgres, anon, authenticated, service_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA supabase_functions GRANT ALL ON SEQUENCES TO postgres, anon, authenticated, service_role;

-- Create supabase_functions tables
CREATE TABLE supabase_functions.migrations (
  version text PRIMARY KEY,
  inserted_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE supabase_functions.hooks (
  id bigserial PRIMARY KEY,
  hook_table_id integer NOT NULL,
  hook_name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  request_id bigint
);

CREATE INDEX supabase_functions_hooks_request_id_idx ON supabase_functions.hooks USING btree (request_id);
CREATE INDEX supabase_functions_hooks_h_table_id_h_name_idx ON supabase_functions.hooks USING btree (hook_table_id, hook_name);

-- Create supabase_functions_admin role
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_roles
    WHERE rolname = 'supabase_functions_admin'
  )
  THEN
    CREATE USER supabase_functions_admin NOINHERIT CREATEROLE LOGIN NOREPLICATION;
  END IF;
END
$$;

-- Grant privileges to supabase_functions_admin
GRANT ALL PRIVILEGES ON SCHEMA supabase_functions TO supabase_functions_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA supabase_functions TO supabase_functions_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA supabase_functions TO supabase_functions_admin;
ALTER USER supabase_functions_admin SET search_path = "supabase_functions";
ALTER table "supabase_functions".migrations OWNER TO supabase_functions_admin;
ALTER table "supabase_functions".hooks OWNER TO supabase_functions_admin;

-- Create http_request function
CREATE FUNCTION supabase_functions.http_request()
  RETURNS trigger
  LANGUAGE plpgsql
  AS $function$
  DECLARE
    request_id bigint;
    payload jsonb;
    url text := TG_ARGV[0]::text;
    method text := TG_ARGV[1]::text;
    headers jsonb DEFAULT '{}'::jsonb;
    params jsonb DEFAULT '{}'::jsonb;
    timeout_ms integer DEFAULT 1000;
  BEGIN
    IF url IS NULL OR url = 'null' THEN
      RAISE EXCEPTION 'url argument is missing';
    END IF;

    IF method IS NULL OR method = 'null' THEN
      RAISE EXCEPTION 'method argument is missing';
    END IF;

    IF TG_ARGV[2] IS NULL OR TG_ARGV[2] = 'null' THEN
      headers = '{"Content-Type": "application/json"}'::jsonb;
    ELSE
      headers = TG_ARGV[2]::jsonb;
    END IF;

    IF TG_ARGV[3] IS NULL OR TG_ARGV[3] = 'null' THEN
      params = '{}'::jsonb;
    ELSE
      params = TG_ARGV[3]::jsonb;
    END IF;

    IF TG_ARGV[4] IS NULL OR TG_ARGV[4] = 'null' THEN
      timeout_ms = 1000;
    ELSE
      timeout_ms = TG_ARGV[4]::integer;
    END IF;

    CASE
      WHEN method = 'GET' THEN
        SELECT http_get INTO request_id FROM net.http_get(
          url,
          params,
          headers,
          timeout_ms
        );
      WHEN method = 'POST' THEN
        payload = jsonb_build_object(
          'old_record', OLD,
          'record', NEW,
          'type', TG_OP,
          'table', TG_TABLE_NAME,
          'schema', TG_TABLE_SCHEMA
        );

        SELECT http_post INTO request_id FROM net.http_post(
          url,
          payload,
          params,
          headers,
          timeout_ms
        );
      ELSE
        RAISE EXCEPTION 'method argument % is invalid', method;
    END CASE;

    INSERT INTO supabase_functions.hooks
      (hook_table_id, hook_name, request_id)
    VALUES
      (TG_RELID, TG_NAME, request_id);

    RETURN NEW;
  END
$function$;

-- Set up pg_net access
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_extension
    WHERE extname = 'pg_net'
  )
  THEN
    GRANT USAGE ON SCHEMA net TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
    REVOKE ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
    REVOKE ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
    GRANT EXECUTE ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    GRANT EXECUTE ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
  END IF;
END
$$;

-- Create event trigger for pg_net
CREATE OR REPLACE FUNCTION extensions.grant_pg_net_access()
RETURNS event_trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_event_trigger_ddl_commands() AS ev
    JOIN pg_extension AS ext
    ON ev.objid = ext.oid
    WHERE ext.extname = 'pg_net'
  )
  THEN
    GRANT USAGE ON SCHEMA net TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
    REVOKE ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
    REVOKE ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
    GRANT EXECUTE ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    GRANT EXECUTE ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
  END IF;
END;
$$;

-- Create event trigger
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_event_trigger
    WHERE evtname = 'issue_pg_net_access'
  ) THEN
    CREATE EVENT TRIGGER issue_pg_net_access ON ddl_command_end WHEN TAG IN ('CREATE EXTENSION')
    EXECUTE PROCEDURE extensions.grant_pg_net_access();
  END IF;
END
$$;

-- Initial migrations
INSERT INTO supabase_functions.migrations (version) VALUES ('initial');
INSERT INTO supabase_functions.migrations (version) VALUES ('20210809183423_update_grants');

-- Set up security for http_request function
ALTER function supabase_functions.http_request() SECURITY DEFINER;
ALTER function supabase_functions.http_request() SET search_path = supabase_functions;
REVOKE ALL ON FUNCTION supabase_functions.http_request() FROM PUBLIC;
GRANT EXECUTE ON FUNCTION supabase_functions.http_request() TO postgres, anon, authenticated, service_role;

-- Grant supabase_functions_admin to postgres
GRANT supabase_functions_admin TO postgres;

-- Remove unused supabase_pg_net_admin role
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_roles
    WHERE rolname = 'supabase_pg_net_admin'
  )
  THEN
    REASSIGN OWNED BY supabase_pg_net_admin TO postgres;
    DROP OWNED BY supabase_pg_net_admin;
    DROP ROLE supabase_pg_net_admin;
  END IF;
END
$$;
`,
				},
			},
			{
				ID:           2,
				Name:         "kong",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("kong:2.8.1"),
				HostInputIDs: []int{1},
				IsPublic:     true,
				Ports: []schema.PortSpec{
					{
						Port:     8000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/",
					Port:                      utils.ToPtr(int32(8000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariablesMounts: []*schema.VariableMount{
					{
						Name: "kong.yml",
						Path: "/kong.yml",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "KONG_DECLARATIVE_CONFIG",
						Value: "/kong.yml",
					},
					{
						Name: "GENERATED_JWT_VARIABLES",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeJWT,
							JWTParams: &schema.JWTParams{
								Issuer:           "supabase",
								SecretOutputKey:  "JWT_SECRET",
								AnonOutputKey:    "SUPABASE_ANON_KEY",
								ServiceOutputKey: "SUPABASE_SERVICE_KEY",
							},
						},
					},
					{
						Name:  "DASHBOARD_USERNAME",
						Value: "admin",
					},
					{
						Name: "DASHBOARD_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
					{
						Name: "kong.yml",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeStringReplace,
						},
						Value: `_format_version: '2.1'
_transform: true

consumers:
  - username: DASHBOARD
  - username: anon
    keyauth_credentials:
      - key: ${SERVICE_2_SUPABASE_ANON_KEY}
  - username: service_role
    keyauth_credentials:
      - key: ${SERVICE_2_SUPABASE_SERVICE_KEY}

acls:
  - consumer: anon
    group: anon
  - consumer: service_role
    group: admin

basicauth_credentials:
- consumer: DASHBOARD
  username: ${SERVICE_2_DASHBOARD_USERNAME}
  password: ${SERVICE_2_DASHBOARD_PASSWORD}

services:
  - name: auth-v1-open
    url: http://${SERVICE_9_KUBE_NAME}.${NAMESPACE}:9999/verify
    routes:
      - name: auth-v1-open
        strip_path: true
        paths:
          - /auth/v1/verify
    plugins:
      - name: cors
  - name: auth-v1
    url: http://${SERVICE_9_KUBE_NAME}.${NAMESPACE}:9999/
    routes:
      - name: auth-v1-all
        strip_path: true
        paths:
          - /auth/v1/
    plugins:
      - name: cors
      - name: key-auth
        config:
          hide_credentials: false
      - name: acl
        config:
          hide_groups_header: true
          allow:
            - admin
            - anon
  - name: rest-v1
    url: http://${SERVICE_8_KUBE_NAME}.${NAMESPACE}:3000/
    routes:
      - name: rest-v1-all
        strip_path: true
        paths:
          - /rest/v1/
    plugins:
      - name: cors
      - name: key-auth
        config:
          hide_credentials: true
      - name: acl
        config:
          hide_groups_header: true
          allow:
            - admin
            - anon
  - name: storage-v1
    url: http://${SERVICE_6_KUBE_NAME}.${NAMESPACE}:5000/
    routes:
      - name: storage-v1-all
        strip_path: true
        paths:
          - /storage/v1/
    plugins:
      - name: cors
  - name: functions-v1
    url: http://${SERVICE_11_KUBE_NAME}.${NAMESPACE}:9000/
    routes:
      - name: functions-v1-all
        strip_path: true
        paths:
          - /functions/v1/
    plugins:
      - name: cors
  - name: analytics-v1
    url: http://${SERVICE_4_KUBE_NAME}.${NAMESPACE}:4000/
    routes:
      - name: analytics-v1-all
        strip_path: true
        paths:
          - /analytics/v1/
  - name: meta
    url: http://${SERVICE_10_KUBE_NAME}.${NAMESPACE}:8080/
    routes:
      - name: meta-all
        strip_path: true
        paths:
          - /pg/
    plugins:
      - name: key-auth
        config:
          hide_credentials: false
      - name: acl
        config:
          hide_groups_header: true
          allow:
            - admin
  - name: dashboard
    url: http://${SERVICE_3_KUBE_NAME}.${NAMESPACE}:3000/
    routes:
      - name: dashboard-all
        strip_path: true
        paths:
          - /
    plugins:
      - name: cors
      - name: basic-auth
        config:
          hide_credentials: true`,
					},
				},
			},
			{
				ID:           3,
				Name:         "studio",
				Type:         schema.ServiceTypeDockerimage,
				Builder:      schema.ServiceBuilderDocker,
				Image:        utils.ToPtr("supabase/studio:2025.04.21-sha-173cc56"),
				DependsOn:    []int{1, 2},
				HostInputIDs: []int{2},
				Ports: []schema.PortSpec{
					{
						Port:     3000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				IsPublic: true,
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/api/profile",
					Port:                      utils.ToPtr(int32(3000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   2,
						TargetName: "SUPABASE_URL",
						IsHost:     true,
					},
					{
						SourceID:   2,
						SourceName: "SUPABASE_ANON_KEY",
						TargetName: "SUPABASE_ANON_KEY",
					},
					{
						SourceID:   2,
						SourceName: "SUPABASE_SERVICE_KEY",
						TargetName: "SUPABASE_SERVICE_KEY",
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "AUTH_JWT_SECRET",
					},
					{
						SourceID:   4,
						SourceName: "LOGFLARE_API_KEY",
						TargetName: "LOGFLARE_API_KEY",
					},
					{
						SourceID:   4,
						TargetName: "LOGFLARE_URL",
						IsHost:     true,
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "SUPABASE_PUBLIC_URL",
						Generator: &schema.ValueGenerator{
							Type:      schema.GeneratorTypeInput,
							InputID:   1,
							AddPrefix: "https://",
						},
					},
				},
			},
			{
				ID:        4,
				Name:      "analytics",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("supabase/logflare:1.12.0"),
				DependsOn: []int{1},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/health",
					Port:                      utils.ToPtr(int32(4000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   10,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				Ports: []schema.PortSpec{
					{
						Port:     4000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "DB_PASSWORD",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "DB_HOSTNAME",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PORT",
						TargetName: "DB_PORT",
					},
					{
						SourceID:       1,
						SourceName:     "DATABASE_PASSWORD",
						TargetName:     "POSTGRES_BACKEND_URL",
						TemplateString: `postgresql://postgres:${DATABASE_PASSWORD}@db:5432/_supabase?sslmode=disable`,
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "DB_USERNAME",
						Value: "postgres",
					},
					{
						Name:  "DB_DATABASE",
						Value: "_supabase",
					},
					{
						Name:  "DB_SCHEMA",
						Value: "_analytics",
					},
					{
						Name: "LOGFLARE_API_KEY",
						Generator: &schema.ValueGenerator{
							Type:     schema.GeneratorTypePassword,
							HashType: utils.ToPtr(schema.ValueHashTypeSHA256),
						},
					},
					{
						Name:  "LOGFLARE_SINGLE_TENANT",
						Value: "true",
					},
					{
						Name:  "LOGFLARE_SINGLE_TENANT_MODE",
						Value: "true",
					},
					{
						Name:  "LOGFLARE_SUPABASE_MODE",
						Value: "true",
					},
					{
						Name:  "LOGFLARE_MIN_CLUSTER_SIZE",
						Value: "1",
					},
					{
						Name:  "POSTGRES_BACKEND_SCHEMA",
						Value: "_analytics",
					},
					{
						Name:  "LOGFLARE_FEATURE_FLAG_OVERRIDE",
						Value: "multibackend=true",
					},
				},
			},
			{
				ID:        5,
				Name:      "vector",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("timberio/vector:0.28.1-alpine"),
				DependsOn: []int{4},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/health",
					Port:                      utils.ToPtr(int32(9001)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariablesMounts: []*schema.VariableMount{
					{
						Name: "vector.yml",
						Path: "/etc/vector/vector.yml",
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   4,
						SourceName: "LOGFLARE_API_KEY",
						TargetName: "LOGFLARE_API_KEY",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name: "vector.yml",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypeStringReplace,
						},
						Value: `api:
  enabled: true
  address: 0.0.0.0:9001

sources:
  docker_host:
    type: docker_logs
    exclude_containers:
      - vector

transforms:
  project_logs:
    type: remap
    inputs:
      - docker_host
    source: |-
      .project = "default"
      .event_message = del(.message)
      .appname = del(.container_name)
      del(.container_created_at)
      del(.container_id)
      del(.source_type)
      del(.stream)
      del(.label)
      del(.image)
      del(.host)
      del(.stream)
  router:
    type: route
    inputs:
      - project_logs
    route:
      kong: 'starts_with(string!(.appname), "kong")'
      auth: 'starts_with(string!(.appname), "auth")'
      rest: 'starts_with(string!(.appname), "rest")'
      realtime: 'starts_with(string!(.appname), "realtime")'
      storage: 'starts_with(string!(.appname), "storage")'
      functions: 'starts_with(string!(.appname), "functions")'
      db: 'starts_with(string!(.appname), "db")'

sinks:
  logflare:
    type: http
    inputs:
      - router.*
    encoding:
      codec: json
    method: post
    request:
      retry_max_duration_secs: 10
    uri: http://${SERVICE_4_KUBE_NAME}.${NAMESPACE}:4000/api/logs?source_name=supabase.logs&api_key=${SERVICE_4_LOGFLARE_API_KEY}`,
					},
				},
			},
			{
				ID:        6,
				Name:      "storage",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("supabase/storage-api:v1.22.7"),
				DependsOn: []int{1, 7},
				Ports: []schema.PortSpec{
					{
						Port:     5000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/status",
					Port:                      utils.ToPtr(int32(5000)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   7,
						SourceName: "MINIO_ROOT_USER",
						TargetName: "AWS_ACCESS_KEY_ID",
					},
					{
						SourceID:   7,
						SourceName: "MINIO_ROOT_PASSWORD",
						TargetName: "AWS_SECRET_ACCESS_KEY",
					},
					{
						SourceID:   7,
						TargetName: "STORAGE_S3_ENDPOINT",
						IsHost:     true,
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "AUTH_JWT_SECRET",
					},
					{
						SourceID:       1,
						SourceName:     "DATABASE_PASSWORD",
						TemplateString: "postgresql://supabase_storage_admin:${DATABASE_PASSWORD}@db:5432/postgres?sslmode=disable",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "STORAGE_BACKEND",
						Value: "s3",
					},
					{
						Name:  "STORAGE_S3_BUCKET",
						Value: "stub",
					},
					{
						Name:  "STORAGE_S3_FORCE_PATH_STYLE",
						Value: "true",
					},
					{
						Name:  "STORAGE_S3_REGION",
						Value: "us-east-1",
					},
				},
				Volumes: []schema.TemplateVolume{
					{
						Name: "storage-data",
						Size: schema.TemplateVolumeSize{
							FromInputID: 3,
						},
						MountPath: "/var/lib/storage",
					},
				},
			},
			{
				ID:      7,
				Name:    "minio",
				Type:    schema.ServiceTypeDockerimage,
				Builder: schema.ServiceBuilderDocker,
				Image:   utils.ToPtr("minio/minio"),
				Ports: []schema.PortSpec{
					{
						Port:     9000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
					{
						Port:     9001,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeExec,
					Command:                   "sleep 5 && exit 0",
					PeriodSeconds:             2,
					TimeoutSeconds:            10,
					StartupFailureThreshold:   5,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "MINIO_ROOT_USER",
						Value: "minioadmin",
					},
					{
						Name: "MINIO_ROOT_PASSWORD",
						Generator: &schema.ValueGenerator{
							Type: schema.GeneratorTypePassword,
						},
					},
				},
			},
			{
				ID:        8,
				Name:      "rest",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("postgrest/postgrest:v12.2.11"),
				DependsOn: []int{1},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "PGRST_JWT_SECRET",
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "PGRST_APP_SETTINGS_JWT_SECRET",
						IsHost:     true,
					},
					{
						SourceID:       1,
						SourceName:     "DATABASE_PASSWORD",
						TargetName:     "PGRST_DB_URI",
						TemplateString: "postgresql://postgres:${DATABASE_PASSWORD}@db:5432/postgres?sslmode=disable",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "PGRST_DB_SCHEMAS",
						Value: "public,storage,graphql_public",
					},
					{
						Name:  "PGRST_DB_ANON_ROLE",
						Value: "anon",
					},
					{
						Name:  "PGRST_APP_SETTINGS_JWT_EXP",
						Value: "3600",
					},
				},
			},
			{
				ID:        9,
				Name:      "auth",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("supabase/gotrue:v2.171.0"),
				DependsOn: []int{1},
				Ports: []schema.PortSpec{
					{
						Port:     9999,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeHTTP,
					Path:                      "/health",
					Port:                      utils.ToPtr(int32(9999)),
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   2,
						TargetName: "API_EXTERNAL_URL",
						IsHost:     true,
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "GOTRUE_JWT_SECRET",
					},
					{
						SourceID:       1,
						SourceName:     "DATABASE_PASSWORD",
						TargetName:     "GOTRUE_DB_DATABASE_URL",
						TemplateString: "postgresql://supabase_auth_admin:${DATABASE_PASSWORD}@db:5432/postgres?sslmode=disable",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "GOTRUE_API_HOST",
						Value: "0.0.0.0",
					},
					{
						Name:  "GOTRUE_API_PORT",
						Value: "9999",
					},
					{
						Name:  "GOTRUE_DB_DRIVER",
						Value: "postgres",
					},
					{
						Name: "GOTRUE_SITE_URL",
						Generator: &schema.ValueGenerator{
							Type:    schema.GeneratorTypeInput,
							InputID: 2,
						},
					},
					{
						Name:  "GOTRUE_JWT_EXP",
						Value: "3600",
					},
				},
			},
			{
				ID:        10,
				Name:      "meta",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("supabase/postgres-meta:v0.88.9"),
				DependsOn: []int{1},
				Ports: []schema.PortSpec{
					{
						Port:     8080,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "PG_META_DB_PASSWORD",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_HOST",
						TargetName: "PG_META_DB_HOST",
					},
					{
						SourceID:   1,
						SourceName: "DATABASE_PORT",
						TargetName: "PG_META_DB_PORT",
					},
				},
				Variables: []schema.TemplateVariable{
					{
						Name:  "PG_META_PORT",
						Value: "8080",
					},
					{
						Name:  "PG_META_DB_NAME",
						Value: "postgres",
					},
					{
						Name:  "PG_META_DB_USER",
						Value: "postgres",
					},
				},
			},
			{
				ID:        11,
				Name:      "functions",
				Type:      schema.ServiceTypeDockerimage,
				Builder:   schema.ServiceBuilderDocker,
				Image:     utils.ToPtr("supabase/edge-runtime:v1.67.4"),
				DependsOn: []int{1, 2},
				Ports: []schema.PortSpec{
					{
						Port:     9000,
						Protocol: utils.ToPtr(schema.ProtocolTCP),
					},
				},
				HealthCheck: &schema.HealthCheck{
					Type:                      schema.HealthCheckTypeExec,
					Command:                   "echo 'Edge Functions is healthy'",
					PeriodSeconds:             5,
					TimeoutSeconds:            5,
					StartupFailureThreshold:   3,
					LivenessFailureThreshold:  3,
					ReadinessFailureThreshold: 3,
				},
				VariableReferences: []schema.TemplateVariableReference{
					{
						SourceID:   1,
						SourceName: "DATABASE_PASSWORD",
						TargetName: "POSTGRES_PASSWORD",
					},
					{
						SourceID:   2,
						SourceName: "JWT_SECRET",
						TargetName: "JWT_SECRET",
					},
					{
						SourceID:   2,
						SourceName: "SUPABASE_ANON_KEY",
						TargetName: "SUPABASE_ANON_KEY",
					},
					{
						SourceID:   2,
						SourceName: "SUPABASE_SERVICE_KEY",
						TargetName: "SUPABASE_SERVICE_ROLE_KEY",
					},
					{
						SourceID:   2,
						TargetName: "SUPABASE_URL",
						IsHost:     true,
					},
					{
						SourceID:       1,
						SourceName:     "DATABASE_PASSWORD",
						TargetName:     "SUPABASE_DB_URL",
						TemplateString: "postgresql://postgres:${DATABASE_PASSWORD}@db:5432/postgres?sslmode=disable",
					},
				},
			},
		},
	}
}
