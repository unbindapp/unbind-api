-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP COLUMN "backup_retention_count", DROP COLUMN "backup_schedule";
