-- modify "service_configs" table
ALTER TABLE "service_configs" DROP COLUMN "volume_name", DROP COLUMN "volume_mount_path", ADD COLUMN "volumes" jsonb NULL;
