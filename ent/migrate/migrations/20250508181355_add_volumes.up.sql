-- modify "service_configs" table
ALTER TABLE "service_configs" ADD COLUMN "volume_name" character varying NULL, ADD COLUMN "volume_mount_path" character varying NULL;
