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
					InitDB: `-- Combined Supabase migration file with default passwords set to 'postgres'

-- Set up realtime
-- defaults to empty publication
create publication supabase_realtime;

-- Supabase super admin
-- (postgres user is already a superuser, no need to alter)

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

-- Set up auth roles for the developer
create role anon nologin noinherit;
create role authenticated nologin noinherit; -- "logged in" user: web_user, app_user, etc
create role service_role nologin noinherit bypassrls; -- allow developers to create JWT's that bypass their policies

create user authenticator noinherit password 'postgres';
grant anon to authenticator;
grant authenticated to authenticator;
grant service_role to authenticator;

grant usage on schema public to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on tables to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on functions to postgres, anon, authenticated, service_role;
alter default privileges in schema public grant all on sequences to postgres, anon, authenticated, service_role;

-- Allow Extensions to be used in the API
grant usage on schema extensions to postgres, anon, authenticated, service_role;

-- Set up namespacing
alter user postgres SET search_path TO public, extensions; -- don't include the "auth" schema

-- These are required so that the users receive grants whenever "postgres" creates tables/function
alter default privileges for user postgres in schema public grant all
    on sequences to postgres, anon, authenticated, service_role;
alter default privileges for user postgres in schema public grant all
    on tables to postgres, anon, authenticated, service_role;
alter default privileges for user postgres in schema public grant all
    on functions to postgres, anon, authenticated, service_role;

-- Set short statement/query timeouts for API roles
alter role anon set statement_timeout = '3s';
alter role authenticated set statement_timeout = '8s';

-- Create auth schema
CREATE SCHEMA IF NOT EXISTS auth AUTHORIZATION postgres;

-- auth.users definition
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

-- auth.refresh_tokens definition
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

-- auth.instances definition
CREATE TABLE auth.instances (
    id uuid NOT NULL,
    uuid uuid NULL,
    raw_base_config text NULL,
    created_at timestamptz NULL,
    updated_at timestamptz NULL,
    CONSTRAINT instances_pkey PRIMARY KEY (id)
);
comment on table auth.instances is 'Auth: Manages users across multiple sites.';

-- auth.audit_log_entries definition
CREATE TABLE auth.audit_log_entries (
    instance_id uuid NULL,
    id uuid NOT NULL,
    payload json NULL,
    created_at timestamptz NULL,
    CONSTRAINT audit_log_entries_pkey PRIMARY KEY (id)
);
CREATE INDEX audit_logs_instance_id_idx ON auth.audit_log_entries USING btree (instance_id);
comment on table auth.audit_log_entries is 'Auth: Audit trail for user actions.';

-- auth.schema_migrations definition
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

-- Storage setup
CREATE SCHEMA IF NOT EXISTS storage AUTHORIZATION postgres;

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
    -- @todo return the last part instead of 2
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

-- Search path adjustments
ALTER ROLE postgres SET search_path TO "$user",public,auth,extensions;
ALTER ROLE postgres SET search_path TO "$user",public,extensions;

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

    alter default privileges for user postgres in schema cron grant all
        on sequences to postgres with grant option;
    alter default privileges for user postgres in schema cron grant all
        on tables to postgres with grant option;
    alter default privileges for user postgres in schema cron grant all
        on functions to postgres with grant option;

    grant all privileges on all tables in schema cron to postgres with grant option;

  END IF;

END;
$$;
CREATE EVENT TRIGGER issue_pg_cron_access ON ddl_command_end WHEN TAG in ('CREATE SCHEMA')
EXECUTE PROCEDURE extensions.grant_pg_cron_access();
COMMENT ON FUNCTION extensions.grant_pg_cron_access IS 'Grants access to pg_cron';

-- Event trigger for pg_net
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
COMMENT ON FUNCTION extensions.grant_pg_net_access IS 'Grants access to pg_net';

DO
$$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_event_trigger
    WHERE evtname = 'issue_pg_net_access'
  ) THEN
    CREATE EVENT TRIGGER issue_pg_net_access
    ON ddl_command_end
    WHEN TAG IN ('CREATE EXTENSION')
    EXECUTE PROCEDURE extensions.grant_pg_net_access();
  END IF;
END
$$;

-- Supabase dashboard user
CREATE ROLE dashboard_user NOSUPERUSER CREATEDB CREATEROLE REPLICATION PASSWORD 'postgres';
GRANT ALL ON DATABASE postgres TO dashboard_user;
GRANT ALL ON SCHEMA auth TO dashboard_user;
GRANT ALL ON SCHEMA extensions TO dashboard_user;
GRANT ALL ON SCHEMA storage TO dashboard_user;
GRANT ALL ON ALL TABLES IN SCHEMA auth TO dashboard_user;
GRANT ALL ON ALL TABLES IN SCHEMA extensions TO dashboard_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA auth TO dashboard_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA storage TO dashboard_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA extensions TO dashboard_user;
GRANT ALL ON ALL ROUTINES IN SCHEMA auth TO dashboard_user;
GRANT ALL ON ALL ROUTINES IN SCHEMA storage TO dashboard_user;
GRANT ALL ON ALL ROUTINES IN SCHEMA extensions TO dashboard_user;

-- Update auth schema permissions
GRANT ALL PRIVILEGES ON SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO supabase_auth_admin;

ALTER table IF EXISTS "auth".users OWNER TO supabase_auth_admin;
ALTER table IF EXISTS "auth".refresh_tokens OWNER TO supabase_auth_admin;
ALTER table IF EXISTS "auth".audit_log_entries OWNER TO supabase_auth_admin;
ALTER table IF EXISTS "auth".instances OWNER TO supabase_auth_admin;
ALTER table IF EXISTS "auth".schema_migrations OWNER TO supabase_auth_admin;

GRANT USAGE ON SCHEMA auth TO postgres;
GRANT ALL ON ALL TABLES IN SCHEMA auth TO postgres, dashboard_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA auth TO postgres, dashboard_user;
GRANT ALL ON ALL ROUTINES IN SCHEMA auth TO postgres, dashboard_user;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_auth_admin IN SCHEMA auth GRANT ALL ON TABLES TO postgres, dashboard_user;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_auth_admin IN SCHEMA auth GRANT ALL ON SEQUENCES TO postgres, dashboard_user;
ALTER DEFAULT PRIVILEGES FOR ROLE supabase_auth_admin IN SCHEMA auth GRANT ALL ON ROUTINES TO postgres, dashboard_user;

-- Create and update realtime schema
CREATE SCHEMA IF NOT EXISTS realtime;
GRANT USAGE ON SCHEMA realtime TO postgres;
GRANT ALL ON ALL TABLES IN SCHEMA realtime TO postgres, dashboard_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA realtime TO postgres, dashboard_user;
GRANT ALL ON ALL ROUTINES IN SCHEMA realtime TO postgres, dashboard_user;

-- Update owner for auth functions
DO $$
BEGIN
    ALTER FUNCTION auth.uid owner to supabase_auth_admin;
EXCEPTION WHEN OTHERS THEN
    RAISE WARNING 'Error encountered when changing owner of auth.uid to supabase_auth_admin';
END $$;

DO $$
BEGIN
    ALTER FUNCTION auth.role owner to supabase_auth_admin;
EXCEPTION WHEN OTHERS THEN
    RAISE WARNING 'Error encountered when changing owner of auth.role to supabase_auth_admin';
END $$;

DO $$
BEGIN
    ALTER FUNCTION auth.email owner to supabase_auth_admin;
EXCEPTION WHEN OTHERS THEN
    RAISE WARNING 'Error encountered when changing owner of auth.email to supabase_auth_admin';
END $$;

-- Update future objects' permissions
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA realtime GRANT ALL ON TABLES TO postgres, dashboard_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA realtime GRANT ALL ON SEQUENCES TO postgres, dashboard_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA realtime GRANT ALL ON ROUTINES TO postgres, dashboard_user;

-- Safe update and session settings
ALTER ROLE authenticator SET session_preload_libraries = 'safeupdate';
ALTER ROLE authenticator set lock_timeout to '8s';

-- PostgreSQL reload schema trigger
CREATE OR REPLACE FUNCTION extensions.pgrst_ddl_watch() RETURNS event_trigger AS $$
DECLARE
  cmd record;
BEGIN
  FOR cmd IN SELECT * FROM pg_event_trigger_ddl_commands()
  LOOP
    IF cmd.command_tag IN (
      'CREATE SCHEMA', 'ALTER SCHEMA'
    , 'CREATE TABLE', 'CREATE TABLE AS', 'SELECT INTO', 'ALTER TABLE'
    , 'CREATE FOREIGN TABLE', 'ALTER FOREIGN TABLE'
    , 'CREATE VIEW', 'ALTER VIEW'
    , 'CREATE MATERIALIZED VIEW', 'ALTER MATERIALIZED VIEW'
    , 'CREATE FUNCTION', 'ALTER FUNCTION'
    , 'CREATE TRIGGER'
    , 'CREATE TYPE', 'ALTER TYPE'
    , 'CREATE RULE'
    , 'COMMENT'
    )
    -- don't notify in case of CREATE TEMP table or other objects created on pg_temp
    AND cmd.schema_name is distinct from 'pg_temp'
    THEN
      NOTIFY pgrst, 'reload schema';
    END IF;
  END LOOP;
END; $$ LANGUAGE plpgsql;

-- Watch drop
CREATE OR REPLACE FUNCTION extensions.pgrst_drop_watch() RETURNS event_trigger AS $$
DECLARE
  obj record;
BEGIN
  FOR obj IN SELECT * FROM pg_event_trigger_dropped_objects()
  LOOP
    IF obj.object_type IN (
      'schema'
    , 'table'
    , 'foreign table'
    , 'view'
    , 'materialized view'
    , 'function'
    , 'trigger'
    , 'type'
    , 'rule'
    )
    AND obj.is_temporary IS false -- no pg_temp objects
    THEN
      NOTIFY pgrst, 'reload schema';
    END IF;
  END LOOP;
END; $$ LANGUAGE plpgsql;

DROP EVENT TRIGGER IF EXISTS pgrst_ddl_watch;
CREATE EVENT TRIGGER pgrst_ddl_watch
  ON ddl_command_end
  EXECUTE PROCEDURE extensions.pgrst_ddl_watch();

DROP EVENT TRIGGER IF EXISTS pgrst_drop_watch;
CREATE EVENT TRIGGER pgrst_drop_watch
  ON sql_drop
  EXECUTE PROCEDURE extensions.pgrst_drop_watch();

-- Set up supautils if available
DO $$
DECLARE
  supautils_exists boolean;
BEGIN
  supautils_exists = (
      select count(*) = 1
      from pg_available_extensions
      where name = 'supautils'
  );

  IF supautils_exists
  THEN
  ALTER ROLE authenticator SET session_preload_libraries = supautils, safeupdate;
  END IF;
END $$;

-- GraphQL setup
create schema if not exists graphql_public;

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

alter default privileges for user postgres in schema graphql_public grant all
    on sequences to postgres, anon, authenticated, service_role;
alter default privileges for user postgres in schema graphql_public grant all
    on tables to postgres, anon, authenticated, service_role;
alter default privileges for user postgres in schema graphql_public grant all
    on functions to postgres, anon, authenticated, service_role;

-- Trigger upon enabling pg_graphql
CREATE OR REPLACE FUNCTION extensions.grant_pg_graphql_access()
RETURNS event_trigger
LANGUAGE plpgsql
AS $func$
DECLARE
    func_is_graphql_resolve bool;
BEGIN
    func_is_graphql_resolve = (
        SELECT n.proname = 'resolve'
        FROM pg_event_trigger_ddl_commands() AS ev
        LEFT JOIN pg_catalog.pg_proc AS n
        ON ev.objid = n.oid
    );

    IF func_is_graphql_resolve
    THEN
        -- Update public wrapper to pass all arguments through to the pg_graphql resolve func
        DROP FUNCTION IF EXISTS graphql_public.graphql;
        create or replace function graphql_public.graphql(
            "operationName" text default null,
            query text default null,
            variables jsonb default null,
            extensions jsonb default null
        )
            returns jsonb
            language sql
        as $$
            select graphql.resolve(
                query := query,
                variables := coalesce(variables, '{}'),
                "operationName" := "operationName",
                extensions := extensions
            );
        $$;

        -- This hook executes when graphql.resolve is created. That is not necessarily the last
        -- function in the extension so we need to grant permissions on existing entities AND
        -- update default permissions to any others that are created after graphql.resolve
        grant usage on schema graphql to postgres, anon, authenticated, service_role;
        grant select on all tables in schema graphql to postgres, anon, authenticated, service_role;
        grant execute on all functions in schema graphql to postgres, anon, authenticated, service_role;
        grant all on all sequences in schema graphql to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on tables to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on functions to postgres, anon, authenticated, service_role;
        alter default privileges in schema graphql grant all on sequences to postgres, anon, authenticated, service_role;

        -- Allow postgres role to allow granting usage on graphql and graphql_public schemas to custom roles
        grant usage on schema graphql_public to postgres with grant option;
        grant usage on schema graphql to postgres with grant option;
    END IF;
END;
$func$;

-- Trigger upon dropping the pg_graphql extension
CREATE OR REPLACE FUNCTION extensions.set_graphql_placeholder()
RETURNS event_trigger
LANGUAGE plpgsql
AS $func$
    DECLARE
    graphql_is_dropped bool;
    BEGIN
    graphql_is_dropped = (
        SELECT ev.schema_name = 'graphql_public'
        FROM pg_event_trigger_dropped_objects() AS ev
        WHERE ev.schema_name = 'graphql_public'
    );

    IF graphql_is_dropped
    THEN
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
    END IF;

    END;
$func$;

DROP EVENT TRIGGER IF EXISTS issue_graphql_placeholder;
CREATE EVENT TRIGGER issue_graphql_placeholder ON sql_drop WHEN TAG in ('DROP EXTENSION')
EXECUTE PROCEDURE extensions.set_graphql_placeholder();
COMMENT ON FUNCTION extensions.set_graphql_placeholder IS 'Reintroduces placeholder function for graphql_public.graphql';

DROP EVENT TRIGGER IF EXISTS issue_pg_graphql_access;
CREATE EVENT TRIGGER issue_pg_graphql_access ON ddl_command_end WHEN TAG in ('CREATE FUNCTION')
EXECUTE PROCEDURE extensions.grant_pg_graphql_access();
COMMENT ON FUNCTION extensions.grant_pg_graphql_access IS 'Grants access to pg_graphql';

-- Auth admin idle timeout
ALTER ROLE supabase_auth_admin SET idle_in_transaction_session_timeout TO 60000;

-- Install pg_graphql if available
DROP EXTENSION IF EXISTS pg_graphql;
DO $$
DECLARE
  graphql_exists boolean;
BEGIN
  graphql_exists = (
      select count(*) = 1
      from pg_available_extensions
      where name = 'pg_graphql'
  );

IF graphql_exists
  THEN
  create extension if not exists pg_graphql;
  END IF;
END $$;

-- Grant postgres role access to manage auth/storage
grant supabase_auth_admin, supabase_storage_admin to postgres;

-- Additional permissions for pg_cron if installed
DO $$
DECLARE
  pg_cron_installed boolean;
BEGIN
  -- checks if pg_cron is enabled   
  pg_cron_installed = (
    select count(*) = 1 
    from pg_available_extensions 
    where name = 'pg_cron'
    and installed_version is not null
  );

  IF pg_cron_installed
  THEN
    grant usage on schema cron to postgres with grant option;
    grant all on all functions in schema cron to postgres with grant option;

    alter default privileges in schema cron grant all on tables to postgres with grant option;
    alter default privileges in schema cron grant all on functions to postgres with grant option;
    alter default privileges in schema cron grant all on sequences to postgres with grant option;

    alter default privileges for user postgres in schema cron grant all
        on sequences to postgres with grant option;
    alter default privileges for user postgres in schema cron grant all
        on tables to postgres with grant option;
    alter default privileges for user postgres in schema cron grant all
        on functions to postgres with grant option;

    grant all privileges on all tables in schema cron to postgres with grant option;

    -- Revoke job access and grant only select
    revoke all on table cron.job from postgres;
    grant select on table cron.job to postgres with grant option;
  END IF;
END $$;

-- pg_net permissions if installed
DO $$
DECLARE
  pg_net_installed boolean;
BEGIN
  -- checks if pg_net is enabled
  pg_net_installed = (
    select count(*) = 1 
    from pg_available_extensions 
    where name = 'pg_net'
    and installed_version is not null
  );

  IF pg_net_installed 
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

    -- For pg_net security
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
    ELSE
      -- For newer versions
      ALTER function net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY INVOKER;
      ALTER function net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) SECURITY INVOKER;

      GRANT ALL ON FUNCTION net.http_get(url text, params jsonb, headers jsonb, timeout_milliseconds integer) TO PUBLIC;
      GRANT ALL ON FUNCTION net.http_post(url text, body jsonb, params jsonb, headers jsonb, timeout_milliseconds integer) TO PUBLIC;
    END IF;
  END IF;
END $$;

-- Role settings and inheritance
ALTER ROLE authenticated inherit;
ALTER ROLE anon inherit;
ALTER ROLE service_role inherit;

-- Grant postgres additional privileges
grant anon, authenticated, service_role to postgres;
grant pg_monitor to postgres;
grant pg_read_all_data, pg_signal_backend to postgres;

-- pgsodium and vault setup
DO $$
DECLARE
  pgsodium_exists boolean;
  vault_exists boolean;
BEGIN
  IF EXISTS (SELECT FROM pg_available_extensions WHERE name = 'supabase_vault' AND default_version != '0.2.8') THEN
    CREATE EXTENSION IF NOT EXISTS supabase_vault;

    -- for some reason extension custom scripts aren't run during AMI build, so
    -- we manually run it here
    grant usage on schema vault to postgres with grant option;
    grant select, delete, truncate, references on vault.secrets, vault.decrypted_secrets to postgres with grant option;
    grant execute on function vault.create_secret, vault.update_secret, vault._crypto_aead_det_decrypt to postgres with grant option;
    grant usage on schema vault to service_role;
    grant select, delete on vault.secrets, vault.decrypted_secrets to service_role;
    grant execute on function vault.create_secret, vault.update_secret, vault._crypto_aead_det_decrypt to service_role;
  ELSE
    pgsodium_exists = (
      select count(*) = 1 
      from pg_available_extensions 
      where name = 'pgsodium'
      and default_version in ('3.1.6', '3.1.7', '3.1.8', '3.1.9')
    );
    
    vault_exists = (
        select count(*) = 1 
        from pg_available_extensions 
        where name = 'supabase_vault'
    );
  
    IF pgsodium_exists 
    THEN
      create extension if not exists pgsodium;
  
      grant pgsodium_keyiduser to postgres with admin option;
      grant pgsodium_keyholder to postgres with admin option;
      grant pgsodium_keymaker to postgres with admin option;
  
      grant execute on function pgsodium.crypto_aead_det_decrypt(bytea, bytea, uuid, bytea) to service_role;
      grant execute on function pgsodium.crypto_aead_det_encrypt(bytea, bytea, uuid, bytea) to service_role;
      grant execute on function pgsodium.crypto_aead_det_keygen to service_role;
  
      IF vault_exists
      THEN
        create extension if not exists supabase_vault;
      END IF;
    END IF;
  END IF;
  
  -- Service role access to pgsodium if available
  IF EXISTS (SELECT FROM pg_roles WHERE rolname = 'pgsodium_keyholder') THEN
    GRANT pgsodium_keyholder to service_role;
  END IF;
END $$;

-- Grant storage admin access to authenticator
grant authenticator to supabase_storage_admin;
revoke anon, authenticated, service_role from supabase_storage_admin;

-- Mask role definition for pgsodium
DO $$
BEGIN
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

-- Set up the PgBouncer auth function
create or replace function pgbouncer.get_auth(p_usename text) returns table (username text, password text)
    language plpgsql security definer
    as $$
begin
    raise debug 'PgBouncer auth request: %', p_usename;

    return query
    select 
        rolname::text, 
        case when rolvaliduntil < now() 
            then null 
            else rolpassword::text 
        end 
    from pg_authid 
    where rolname=$1 and rolcanlogin;
end;
$$;

alter function pgbouncer.get_auth owner to supabase_auth_admin;
grant execute on function pgbouncer.get_auth(p_usename text) to postgres;

-- Orioledb extension if available
do $$ 
begin 
    if exists (select 1 from pg_available_extensions where name = 'orioledb') then
        if not exists (select 1 from pg_extension where extname = 'orioledb') then
            create extension if not exists orioledb;
        end if;
    end if;
end $$;

-- Move orioledb to extensions schema if installed in public
do $$
declare
    ext_schema text;
    extensions_schema_exists boolean;
begin
    -- check if the "extensions" schema exists
    select exists (
        select 1 from pg_namespace where nspname = 'extensions'
    ) into extensions_schema_exists;

    if extensions_schema_exists then
        -- check if the "orioledb" extension is in the "public" schema
        select nspname into ext_schema
        from pg_extension e
        join pg_namespace n on e.extnamespace = n.oid
        where extname = 'orioledb';

        if ext_schema = 'public' then
            execute 'alter extension orioledb set schema extensions';
        end if;
    end if;
end $$;

-- PGMQ meta data setup
do $$
begin
    -- Check if the pgmq.meta table exists
    if exists (
        select
            1
        from
            pg_catalog.pg_class c
        join pg_catalog.pg_namespace n
            on c.relnamespace = n.oid
        where
            n.nspname = 'pgmq'
            and c.relname = 'meta'
            and c.relkind = 'r' -- regular table
            -- Make sure only expected columns exist and are correctly named
            and (
                select array_agg(attname::text order by attname)
                from pg_catalog.pg_attribute a
                where
                a.attnum > 0
                and a.attrelid = c.oid 
            ) = array['created_at', 'is_partitioned', 'is_unlogged', 'queue_name']::text[]
    ) then
        -- Insert data into pgmq.meta for all tables matching the naming pattern 'pgmq.q_<queue_name>'
        insert into pgmq.meta (queue_name, is_partitioned, is_unlogged, created_at)
        select
            substring(c.relname from 3) as queue_name,
            false as is_partitioned,
            case when c.relpersistence = 'u' then true else false end as is_unlogged,
            now() as created_at
        from
            pg_catalog.pg_class c
            join pg_catalog.pg_namespace n
                on c.relnamespace = n.oid
        where
            n.nspname = 'pgmq'
            and c.relname like 'q_%'
            and c.relkind in ('r', 'p', 'u')
        on conflict (queue_name) do nothing;
    end if;
end $$;

-- Reattach any detached PGMQ tables
do $$
declare
    ext_exists boolean;
    tbl record;
begin
    -- check if pgmq extension is installed
    select exists(select 1 from pg_extension where extname = 'pgmq') into ext_exists;

    if ext_exists then
        for tbl in
            select c.relname as table_name
            from pg_class c
            join pg_namespace n on c.relnamespace = n.oid
            where n.nspname = 'pgmq'
              and c.relkind in ('r', 'u')  -- include ordinary and unlogged tables
              and (c.relname like 'q\_%' or c.relname like 'a\_%')
              and c.oid not in (
                  select d.objid
                  from pg_depend d
                  join pg_extension e on d.refobjid = e.oid
                  where e.extname = 'pgmq'
                    and d.classid = 'pg_class'::regclass
                    and d.deptype = 'e'
              )
        loop
            execute format('alter extension pgmq add table pgmq.%I', tbl.table_name);
        end loop;
    end if;
end;
$$;

-- Disable logging for admin roles
alter role postgres set log_statement = none;
alter role supabase_auth_admin set log_statement = none;
alter role supabase_storage_admin set log_statement = none;

-- Remove storage and auth migrations from standard roles
revoke supabase_storage_admin from postgres;
revoke create on schema storage from postgres;
revoke all on storage.migrations from anon, authenticated, service_role, postgres;

revoke supabase_auth_admin from postgres;
revoke create on schema auth from postgres;
revoke all on auth.schema_migrations from dashboard_user, postgres;

-- Permissions for extensions
grant all privileges on all tables in schema extensions to postgres with grant option;
grant all privileges on all routines in schema extensions to postgres with grant option;
grant all privileges on all sequences in schema extensions to postgres with grant option;
alter default privileges in schema extensions grant all on tables to postgres with grant option;
alter default privileges in schema extensions grant all on routines to postgres with grant option;
alter default privileges in schema extensions grant all on sequences to postgres with grant option;

-- Security for large objects
alter function pg_catalog.lo_export owner to postgres;
alter function pg_catalog.lo_import(text) owner to postgres;
alter function pg_catalog.lo_import(text, oid) owner to postgres;
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
						TemplateString:            `postgresql://postgres:${DATABASE_PASSWORD}@${DATABASE_HOST}:5432/_supabase?sslmode=disable`,
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
						SourceName:     "DATABASE_HOST",
						TargetName:     "DATABASE_URL",
						TemplateString: "postgresql://supabase_storage_admin:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
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
						SourceName:     "DATABASE_HOST",
						TargetName:     "PGRST_DB_URI",
						TemplateString: "postgresql://authenticator:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
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
						SourceName:                "DATABASE_HOST",
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
						SourceID:       1,
						SourceName:     "DATABASE_HOST",
						TargetName:     "SUPABASE_DB_URL",
						TemplateString: "postgresql://supabase_functions_admin:postgres@${DATABASE_HOST}:5432/postgres?sslmode=disable",
					},
				},
			},
		},
	}
}
