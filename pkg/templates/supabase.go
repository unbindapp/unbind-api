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
-- Set up realtime
create publication supabase_realtime;

-- Supabase super admin
alter user supabase_admin with superuser createdb createrole replication bypassrls;

-- Supabase replication user
create user supabase_replication_admin with login replication password 'postgres';

-- Supabase read-only user
create role supabase_read_only_user with login bypassrls password 'postgres';
grant pg_read_all_data to supabase_read_only_user;

-- Extension namespacing
create schema if not exists extensions;
create extension if not exists "uuid-ossp" with schema extensions;
create extension if not exists pgcrypto with schema extensions;
do $$
begin 
    if exists (select 1 from pg_available_extensions where name = 'pgjwt') then
        if not exists (select 1 from pg_extension where extname = 'pgjwt') then
            if current_setting('server_version_num')::int / 10000 = 15 then
                create extension if not exists pgjwt with schema "extensions" cascade;
            end if;
        end if;
    end if;
end $$;

-- Create required schemas
CREATE SCHEMA IF NOT EXISTS realtime;
CREATE SCHEMA IF NOT EXISTS graphql_public;
CREATE SCHEMA IF NOT EXISTS net;
ALTER SCHEMA net OWNER TO postgres;

-- Set up auth roles for the developer
create role anon nologin noinherit;
create role authenticated nologin noinherit;
create role service_role nologin noinherit bypassrls;

create user authenticator noinherit password 'postgres';
grant anon to authenticator;
grant authenticated to authenticator;
grant service_role to authenticator;
grant supabase_admin to authenticator;

grant usage on schema public to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on tables to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on functions to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on sequences to postgres, anon, authenticated, service_role;

-- Allow Extensions to be used in the API
grant usage on schema extensions to postgres, anon, authenticated, service_role;

-- Set up namespacing
alter user supabase_admin SET search_path TO public, extensions;

-- These are required so that the users receive grants whenever "supabase_admin" creates tables/function
alter default privileges for user supabase_admin in schema public grant all
    on sequences to postgres, anon, authenticated, service_role;
alter default privileges for user supabase_admin in schema public grant all
    on tables to postgres, anon, authenticated, service_role;
alter default privileges for user supabase_admin in schema public grant all
    on functions to postgres, anon, authenticated, service_role;

-- Set short statement/query timeouts for API roles
alter role anon set statement_timeout = '3s';
alter role authenticated set statement_timeout = '8s';

-- Create auth schema and tables
CREATE SCHEMA IF NOT EXISTS auth AUTHORIZATION supabase_admin;

CREATE TABLE auth.users (
    instance_id uuid NULL,
    id uuid NOT NULL UNIQUE,
    aud varchar(255) NULL,
    "role" varchar(255) NULL,
    email varchar(255) NULL UNIQUE,
    encrypted_password varchar(255) NULL,
    confirmed_at timestamptz NULL,
    invited_at timestamptz NULL,
    confirmation_token varchar(255) NULL,
    confirmation_sent_at timestamptz NULL,
    recovery_token varchar(255) NULL,
    recovery_sent_at timestamptz NULL,
    email_change_token varchar(255) NULL,
    email_change varchar(255) NULL,
    email_change_sent_at timestamptz NULL,
    last_sign_in_at timestamptz NULL,
    raw_app_meta_data jsonb NULL,
    raw_user_meta_data jsonb NULL,
    is_super_admin bool NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
CREATE INDEX users_instance_id_email_idx ON auth.users USING btree (instance_id, email);
CREATE INDEX users_instance_id_idx ON auth.users USING btree (instance_id);
comment on table auth.users is 'Auth: Stores user login data within a secure schema.';

CREATE TABLE auth.refresh_tokens (
    instance_id uuid NULL,
    id bigserial NOT NULL,
    "token" varchar(255) NULL,
    user_id varchar(255) NULL,
    revoked bool NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id)
);
CREATE INDEX refresh_tokens_instance_id_idx ON auth.refresh_tokens USING btree (instance_id);
CREATE INDEX refresh_tokens_instance_id_user_id_idx ON auth.refresh_tokens USING btree (instance_id, user_id);
CREATE INDEX refresh_tokens_token_idx ON auth.refresh_tokens USING btree (token);
comment on table auth.refresh_tokens is 'Auth: Store of tokens used to refresh JWT tokens once they expire.';

CREATE TABLE auth.instances (
    id uuid NOT NULL,
    uuid uuid NULL,
    raw_base_config text NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    CONSTRAINT instances_pkey PRIMARY KEY (id)
);
comment on table auth.instances is 'Auth: Manages users across multiple sites.';

CREATE TABLE auth.audit_log_entries (
    instance_id uuid NULL,
    id uuid NOT NULL,
    payload json NULL,
    created_at timestamptz NULL,
    CONSTRAINT audit_log_entries_pkey PRIMARY KEY (id)
);
CREATE INDEX audit_logs_instance_id_idx ON auth.audit_log_entries USING btree (instance_id);
comment on table auth.audit_log_entries is 'Auth: Audit trail for user actions.';

CREATE TABLE auth.schema_migrations (
    "version" varchar(255) NOT NULL,
    CONSTRAINT schema_migrations_pkey PRIMARY KEY ("version")
);
comment on table auth.schema_migrations is 'Auth: Manages updates to the auth system.';

INSERT INTO auth.schema_migrations (version)
VALUES  ('20171026211738'),
        ('20171026211808'),
        ('20171026211834'),
        ('20180103212743'),
        ('20180108183307'),
        ('20180119214651'),
        ('20180125194653');

-- Gets the User ID from the request cookie
create or replace function auth.uid() returns uuid as $$
  select nullif(current_setting('request.jwt.claim.sub', true), '')::uuid;
$$ language sql stable;

-- Gets the User ID from the request cookie
create or replace function auth.role() returns text as $$
  select nullif(current_setting('request.jwt.claim.role', true), '')::text;
$$ language sql stable;

-- Gets the User email
create or replace function auth.email() returns text as $$
  select nullif(current_setting('request.jwt.claim.email', true), '')::text;
$$ language sql stable;

-- usage on auth functions to API roles
GRANT USAGE ON SCHEMA auth TO anon, authenticated, service_role;

-- Supabase auth admin
CREATE USER supabase_auth_admin NOINHERIT CREATEROLE LOGIN NOREPLICATION PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO supabase_auth_admin;
ALTER USER supabase_auth_admin SET search_path = "auth";
ALTER table "auth".users OWNER TO supabase_auth_admin;
ALTER table "auth".refresh_tokens OWNER TO supabase_auth_admin;
ALTER table "auth".audit_log_entries OWNER TO supabase_auth_admin;
ALTER table "auth".instances OWNER TO supabase_auth_admin;
ALTER table "auth".schema_migrations OWNER TO supabase_auth_admin;

-- Create storage schema and tables
CREATE SCHEMA IF NOT EXISTS storage AUTHORIZATION supabase_admin;

grant usage on schema storage to postgres, anon, authenticated, service_role;
alter default privileges in schema storage grant all on tables to postgres, anon, authenticated, service_role;
alter default privileges in schema storage grant all on functions to postgres, anon, authenticated, service_role;
alter default privileges in schema storage grant all on sequences to postgres, anon, authenticated, service_role;

CREATE TABLE "storage"."buckets" (
    "id" text not NULL,
    "name" text NOT NULL,
    "owner" uuid,
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    CONSTRAINT "buckets_owner_fkey" FOREIGN KEY ("owner") REFERENCES "auth"."users"("id"),
    PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "bname" ON "storage"."buckets" USING BTREE ("name");

CREATE TABLE "storage"."objects" (
    "id" uuid NOT NULL DEFAULT extensions.uuid_generate_v4(),
    "bucket_id" text,
    "name" text,
    "owner" uuid,
    "created_at" timestamptz DEFAULT now(),
    "updated_at" timestamptz DEFAULT now(),
    "last_accessed_at" timestamptz DEFAULT now(),
    "metadata" jsonb,
    CONSTRAINT "objects_bucketId_fkey" FOREIGN KEY ("bucket_id") REFERENCES "storage"."buckets"("id"),
    CONSTRAINT "objects_owner_fkey" FOREIGN KEY ("owner") REFERENCES "auth"."users"("id"),
    PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "bucketid_objname" ON "storage"."objects" USING BTREE ("bucket_id","name");
CREATE INDEX name_prefix_search ON storage.objects(name text_pattern_ops);

ALTER TABLE storage.objects ENABLE ROW LEVEL SECURITY;

CREATE FUNCTION storage.foldername(name text)
 RETURNS text[]
 LANGUAGE plpgsql
AS $function$
DECLARE
_parts text[];
BEGIN
    select string_to_array(name, '/') into _parts;
    return _parts[1:array_length(_parts,1)-1];
END
$function$;

CREATE FUNCTION storage.filename(name text)
 RETURNS text
 LANGUAGE plpgsql
AS $function$
DECLARE
_parts text[];
BEGIN
    select string_to_array(name, '/') into _parts;
    return _parts[array_length(_parts,1)];
END
$function$;

CREATE FUNCTION storage.extension(name text)
 RETURNS text
 LANGUAGE plpgsql
AS $function$
DECLARE
_parts text[];
_filename text;
BEGIN
    select string_to_array(name, '/') into _parts;
    select _parts[array_length(_parts,1)] into _filename;
    return split_part(_filename, '.', 2);
END
$function$;

CREATE FUNCTION storage.search(prefix text, bucketname text, limits int DEFAULT 100, levels int DEFAULT 1, offsets int DEFAULT 0)
 RETURNS TABLE (
    name text,
    id uuid,
    updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    last_accessed_at TIMESTAMPTZ,
    metadata jsonb
  )
 LANGUAGE plpgsql
AS $function$
DECLARE
_bucketId text;
BEGIN
    -- will be replaced by migrations when server starts
    -- saving space for cloud-init
END
$function$;

-- create migrations table
CREATE TABLE IF NOT EXISTS storage.migrations (
  id integer PRIMARY KEY,
  name varchar(100) UNIQUE NOT NULL,
  hash varchar(40) NOT NULL,
  executed_at timestamp DEFAULT current_timestamp
);

-- Supabase storage admin
CREATE USER supabase_storage_admin NOINHERIT CREATEROLE LOGIN NOREPLICATION PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON SCHEMA storage TO supabase_storage_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA storage TO supabase_storage_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA storage TO supabase_storage_admin;
ALTER USER supabase_storage_admin SET search_path = "storage";
ALTER table "storage".objects owner to supabase_storage_admin;
ALTER table "storage".buckets owner to supabase_storage_admin;
ALTER table "storage".migrations OWNER TO supabase_storage_admin;
ALTER function "storage".foldername(text) owner to supabase_storage_admin;
ALTER function "storage".filename(text) owner to supabase_storage_admin;
ALTER function "storage".extension(text) owner to supabase_storage_admin;
ALTER function "storage".search(text,text,int,int,int) owner to supabase_storage_admin;

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
    IF NOT EXISTS (
      SELECT 1
      FROM pg_roles
      WHERE rolname = 'supabase_functions_admin'
    )
    THEN
      CREATE USER supabase_functions_admin NOINHERIT CREATEROLE LOGIN NOREPLICATION PASSWORD 'postgres';
    END IF;

    GRANT USAGE ON SCHEMA net TO supabase_functions_admin, postgres, anon, authenticated, service_role;

    IF EXISTS (
      SELECT FROM pg_extension
      WHERE extname = 'pg_net'
      AND extversion IN ('0.2', '0.6', '0.7', '0.7.1', '0.8', '0.10.0', '0.11.0')
    ) THEN
      ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;
      ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY DEFINER;

      ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;
      ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SET search_path = net;

      REVOKE ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;
      REVOKE ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) FROM PUBLIC;

      GRANT EXECUTE ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
      GRANT EXECUTE ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO supabase_functions_admin, postgres, anon, authenticated, service_role;
    END IF;
  END IF;
END;
$$;

CREATE EVENT TRIGGER issue_pg_net_access ON ddl_command_end
  WHEN TAG IN ('CREATE EXTENSION')
  EXECUTE FUNCTION extensions.grant_pg_net_access();

-- GraphQL Placeholder Entrypoint
create or replace function graphql_public.graphql(
    "operationName" text default null,
    query text default null,
    variables jsonb default null,
    extensions jsonb default null
)
    returns jsonb
    language plpgsql
as $$
    DECLARE
        server_version float;
    BEGIN
        server_version = (SELECT (SPLIT_PART((select version()), ' ', 2))::float);

        IF server_version >= 14 THEN
            RETURN jsonb_build_object(
                'errors', jsonb_build_array(
                    jsonb_build_object(
                        'message', 'pg_graphql extension is not enabled.'
                    )
                )
            );
        ELSE
            RETURN jsonb_build_object(
                'errors', jsonb_build_array(
                    jsonb_build_object(
                        'message', 'pg_graphql is only available on projects running Postgres 14 onwards.'
                    )
                )
            );
        END IF;
    END;
$$;

grant usage on schema graphql_public to postgres, anon, authenticated, service_role;
alter default privileges in schema graphql_public grant all on tables to postgres, anon, authenticated, service_role;
alter default privileges in schema graphql_public grant all on functions to postgres, anon, authenticated, service_role;
alter default privileges in schema graphql_public grant all on sequences to postgres, anon, authenticated, service_role;

alter default privileges for user supabase_admin in schema graphql_public grant all
    on sequences to postgres, anon, authenticated, service_role;
alter default privileges for user supabase_admin in schema graphql_public grant all
    on tables to postgres, anon, authenticated, service_role;
alter default privileges for user supabase_admin in schema graphql_public grant all
    on functions to postgres, anon, authenticated, service_role;

-- Trigger for pg_cron
CREATE OR REPLACE FUNCTION extensions.grant_pg_cron_access()
RETURNS event_trigger
LANGUAGE plpgsql
AS $$
DECLARE
  schema_is_cron bool;
BEGIN
  schema_is_cron = (
    SELECT n.nspname = 'cron'
    FROM pg_event_trigger_ddl_commands() AS ev
    LEFT JOIN pg_catalog.pg_namespace AS n
      ON ev.objid = n.oid
  );

  IF schema_is_cron
  THEN
    grant usage on schema cron to postgres with grant option;

    alter default privileges in schema cron grant all on tables to postgres with grant option;
    alter default privileges in schema cron grant all on functions to postgres with grant option;
    alter default privileges in schema cron grant all on sequences to postgres with grant option;

    alter default privileges for user supabase_admin in schema cron grant all
        on sequences to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on tables to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on functions to postgres with grant option;

    grant all privileges on all tables in schema cron to postgres with grant option;
  END IF;
END;
$$;

CREATE EVENT TRIGGER issue_pg_cron_access ON ddl_command_end WHEN TAG in ('CREATE SCHEMA')
EXECUTE PROCEDURE extensions.grant_pg_cron_access();

-- Set up search paths
ALTER ROLE supabase_admin SET search_path TO "$user",public,auth,extensions;
ALTER ROLE postgres SET search_path TO "$user",public,extensions;

-- Set timeouts
ALTER ROLE authenticator SET statement_timeout = '8s';
ALTER ROLE authenticator SET lock_timeout = '8s';
ALTER ROLE supabase_auth_admin SET idle_in_transaction_session_timeout TO 60000;

-- Grant roles
grant supabase_auth_admin, supabase_storage_admin to postgres;
grant anon, authenticated, service_role to postgres;
grant anon, authenticated, service_role to supabase_storage_admin;
grant authenticator to supabase_storage_admin;

-- Set inheritance
ALTER ROLE authenticated inherit;
ALTER ROLE anon inherit;
ALTER ROLE service_role inherit;

-- Revoke unnecessary permissions
revoke supabase_storage_admin from postgres;
revoke create on schema storage from postgres;
revoke all on storage.migrations from anon, authenticated, service_role, postgres;

revoke supabase_auth_admin from postgres;
revoke create on schema auth from postgres;
revoke all on auth.schema_migrations from dashboard_user, postgres;

-- Set logging
alter role supabase_admin set log_statement = none;
alter role supabase_auth_admin set log_statement = none;
alter role supabase_storage_admin set log_statement = none;

-- Set up extensions and their dependencies
DO $$
BEGIN
  -- Set up pg_net
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'pg_net') THEN
    ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY INVOKER;
    ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY INVOKER;

    REVOKE EXECUTE ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) FROM supabase_functions_admin, postgres, anon, authenticated, service_role;
    REVOKE EXECUTE ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) FROM supabase_functions_admin, postgres, anon, authenticated, service_role;

    GRANT ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO PUBLIC;
    GRANT ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO PUBLIC;
  END IF;

  -- Set up pg_cron
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'pg_cron') THEN
    grant usage on schema cron to postgres with grant option;
    grant all on all functions in schema cron to postgres with grant option;

    alter default privileges in schema cron grant all on tables to postgres with grant option;
    alter default privileges in schema cron grant all on functions to postgres with grant option;
    alter default privileges in schema cron grant all on sequences to postgres with grant option;

    alter default privileges for user supabase_admin in schema cron grant all
        on sequences to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on tables to postgres with grant option;
    alter default privileges for user supabase_admin in schema cron grant all
        on functions to postgres with grant option;

    grant all privileges on all tables in schema cron to postgres with grant option;
  END IF;

  -- Set up pg_graphql
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'pg_graphql') THEN
    create extension if not exists pg_graphql;
  END IF;

  -- Set up orioledb
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'orioledb') THEN
    create extension if not exists orioledb;
  END IF;

  -- Set up pgmq
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'pgmq') THEN
    -- Check if the pgmq.meta table exists
    IF EXISTS (
      SELECT 1
      FROM pg_catalog.pg_class c
      JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
      WHERE n.nspname = 'pgmq'
        AND c.relname = 'meta'
        AND c.relkind = 'r'
        AND (
          SELECT array_agg(attname::text ORDER BY attname)
          FROM pg_catalog.pg_attribute a
          WHERE a.attnum > 0
            AND a.attrelid = c.oid 
        ) = array['created_at', 'is_partitioned', 'is_unlogged', 'queue_name']::text[]
    ) THEN
      -- Insert data into pgmq.meta for all tables matching the naming pattern
      INSERT INTO pgmq.meta (queue_name, is_partitioned, is_unlogged, created_at)
      SELECT
        substring(c.relname from 3) as queue_name,
        false as is_partitioned,
        CASE WHEN c.relpersistence = 'u' THEN true ELSE false END as is_unlogged,
        now() as created_at
      FROM pg_catalog.pg_class c
      JOIN pg_catalog.pg_namespace n ON c.relnamespace = n.oid
      WHERE n.nspname = 'pgmq'
        AND c.relname like 'q_%'
        AND c.relkind in ('r', 'p', 'u')
      ON CONFLICT (queue_name) DO NOTHING;

      -- Re-attach queue tables
      FOR tbl IN
        SELECT c.relname as table_name
        FROM pg_class c
        JOIN pg_namespace n ON c.relnamespace = n.oid
        WHERE n.nspname = 'pgmq'
          AND c.relkind in ('r', 'u')
          AND (c.relname like 'q\_%' or c.relname like 'a\_%')
          AND c.oid NOT IN (
            SELECT d.objid
            FROM pg_depend d
            JOIN pg_extension e ON d.refobjid = e.oid
            WHERE e.extname = 'pgmq'
              AND d.classid = 'pg_class'::regclass
              AND d.deptype = 'e'
          )
      LOOP
        EXECUTE format('alter extension pgmq add table pgmq.%I', tbl.table_name);
      END LOOP;
    END IF;
  END IF;

  -- Set up pgsodium
  IF EXISTS (SELECT FROM pg_extension WHERE extname = 'pgsodium') THEN
    CREATE OR REPLACE FUNCTION pgsodium.mask_role(masked_role regrole, source_name text, view_name text)
    RETURNS void
    LANGUAGE plpgsql
    SECURITY DEFINER
    SET search_path TO ''
    AS $function$
    BEGIN
      EXECUTE format(
        'GRANT SELECT ON pgsodium.key TO %s',
        masked_role);

      EXECUTE format(
        'GRANT pgsodium_keyiduser, pgsodium_keyholder TO %s',
        masked_role);

      EXECUTE format(
        'GRANT ALL ON %I TO %s',
        view_name,
        masked_role);
      RETURN;
    END
    $function$;
  END IF;
END $$;

-- Grant additional roles
grant pg_monitor to postgres;
grant pg_read_all_data, pg_signal_backend to postgres;
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
						SourceID:                  1,
						SourceName:                "DATABASE_PASSWORD",
						TargetName:                "POSTGRES_BACKEND_URL",
						AdditionalTemplateSources: []string{"DATABASE_HOST}"},
						TemplateString:            `postgresql://supabase_admin:${DATABASE_PASSWORD}@${DATABASE_HOST}:5432/_supabase?sslmode=disable`,
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
						SourceID:                  1,
						SourceName:                "DATABASE_PASSWORD",
						TargetName:                "DATABASE_URL",
						AdditionalTemplateSources: []string{"DATABASE_HOST"},
						TemplateString:            "postgresql://supabase_storage_admin:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
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
						SourceID:                  1,
						SourceName:                "DATABASE_PASSWORD",
						TargetName:                "PGRST_DB_URI",
						AdditionalTemplateSources: []string{"DATABASE_HOST"},
						TemplateString:            "postgresql://authenticator:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
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
						SourceID:                  1,
						SourceName:                "DATABASE_PASSWORD",
						TargetName:                "GOTRUE_DB_DATABASE_URL",
						AdditionalTemplateSources: []string{"DATABASE_HOST"},
						TemplateString:            "postgresql://supabase_auth_admin:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
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
						SourceID:                  1,
						SourceName:                "DATABASE_PASSWORD",
						TargetName:                "SUPABASE_DB_URL",
						AdditionalTemplateSources: []string{"DATABASE_HOST"},
						TemplateString:            "postgresql://supabase_functions_admin:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
					},
				},
			},
		},
	}
}
