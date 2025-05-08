-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP COLUMN "volume_mount_path", DROP COLUMN "volume_name";
