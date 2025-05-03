-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP CONSTRAINT "service_configs_s3_endpoints_service_backup_endpoint", DROP COLUMN "s3_backup_endpoint_id", ADD COLUMN "s3_backup_source_id" uuid NULL;
-- reverse: create "s3_endpoints" table
DROP TABLE "s3_endpoints";
