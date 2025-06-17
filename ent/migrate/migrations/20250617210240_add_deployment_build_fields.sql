-- +goose Up
-- modify "deployments" table
DO $$
BEGIN
    -- Add columns if they don't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'builder') THEN
        ALTER TABLE "deployments" ADD COLUMN "builder" character varying;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'railpack_builder_install_command') THEN
        ALTER TABLE "deployments" ADD COLUMN "railpack_builder_install_command" character varying;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'railpack_builder_build_command') THEN
        ALTER TABLE "deployments" ADD COLUMN "railpack_builder_build_command" character varying;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'run_command') THEN
        ALTER TABLE "deployments" ADD COLUMN "run_command" character varying;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'docker_builder_dockerfile_path') THEN
        ALTER TABLE "deployments" ADD COLUMN "docker_builder_dockerfile_path" character varying;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'deployments' AND column_name = 'docker_builder_build_context') THEN
        ALTER TABLE "deployments" ADD COLUMN "docker_builder_build_context" character varying;
    END IF;
END $$;

-- Populate the new columns from service_configs
UPDATE deployments 
SET 
    builder = sc.builder,
    railpack_builder_install_command = sc.railpack_builder_install_command,
    railpack_builder_build_command = sc.railpack_builder_build_command,
    run_command = sc.run_command,
    docker_builder_dockerfile_path = sc.docker_builder_dockerfile_path,
    docker_builder_build_context = sc.docker_builder_build_context
FROM service_configs sc
WHERE deployments.service_id = sc.service_id;

-- Now set builder as NOT NULL since it's required in service_configs
ALTER TABLE "deployments" ALTER COLUMN "builder" SET NOT NULL;

-- +goose Down
-- reverse: modify "deployments" table
ALTER TABLE "deployments" 
    DROP COLUMN IF EXISTS "docker_builder_build_context", 
    DROP COLUMN IF EXISTS "docker_builder_dockerfile_path", 
    DROP COLUMN IF EXISTS "run_command", 
    DROP COLUMN IF EXISTS "railpack_builder_build_command", 
    DROP COLUMN IF EXISTS "railpack_builder_install_command", 
    DROP COLUMN IF EXISTS "builder";