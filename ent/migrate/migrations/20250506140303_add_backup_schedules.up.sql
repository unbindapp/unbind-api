-- modify "service_configs" table
ALTER TABLE "service_configs" ADD COLUMN "backup_schedule" character varying NOT NULL DEFAULT '5 5 * * *', ADD COLUMN "backup_retention_count" bigint NOT NULL DEFAULT 3;
