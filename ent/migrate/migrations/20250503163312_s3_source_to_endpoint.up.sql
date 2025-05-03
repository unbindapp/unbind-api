-- Rename the table from "s3_sources" to "s3_endpoints"
ALTER TABLE "s3_sources" RENAME TO "s3_endpoints";

-- Update the constraint name to reflect the new table name
ALTER TABLE "s3_endpoints" RENAME CONSTRAINT "s3_sources_teams_s3_sources" TO "s3_endpoints_teams_s3_endpoints";

-- Rename column in "service_configs" table
ALTER TABLE "service_configs" 
    RENAME COLUMN "s3_backup_source_id" TO "s3_backup_endpoint_id";

-- modify "service_configs" table
ALTER TABLE "service_configs" DROP CONSTRAINT "service_configs_s3_sources_service_backup_source";

-- Add the foreign key constraint 
ALTER TABLE "service_configs"
    ADD CONSTRAINT "service_configs_s3_endpoints_service_backup_endpoint" 
        FOREIGN KEY ("s3_backup_endpoint_id") 
        REFERENCES "s3_endpoints" ("id") 
        ON UPDATE NO ACTION 
        ON DELETE SET NULL;