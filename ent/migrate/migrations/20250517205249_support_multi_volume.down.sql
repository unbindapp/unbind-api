-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP COLUMN "volumes", ADD COLUMN "volume_mount_path" character varying NULL, ADD COLUMN "volume_name" character varying NULL;
