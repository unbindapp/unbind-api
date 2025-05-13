-- Drop the foreign key constraint
ALTER TABLE "service_configs" 
    DROP CONSTRAINT "service_configs_s3_endpoints_service_backup_endpoint";

-- Rename column in "service_configs" table back
ALTER TABLE "service_configs" 
    RENAME COLUMN "s3_backup_endpoint_id" TO "s3_backup_source_id";

-- Rename the constraint back
ALTER TABLE "s3_endpoints" RENAME CONSTRAINT "s3_endpoints_teams_s3_endpoints" TO "s3_sources_teams_s3_sources";

-- Rename the table back from "s3_endpoints" to "s3_sources"
ALTER TABLE "s3_endpoints" RENAME TO "s3_sources";

-- Add the foreign key constraint back
ALTER TABLE "service_configs"
    ADD CONSTRAINT "service_configs_s3_sources_service_backup_source" 
        FOREIGN KEY ("s3_backup_source_id") 
        REFERENCES "s3_sources" ("id") 
        ON UPDATE NO ACTION 
        ON DELETE SET NULL;
