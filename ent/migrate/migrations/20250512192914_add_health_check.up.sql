-- modify "service_configs" table
ALTER TABLE "service_configs" ADD COLUMN "health_check" jsonb NULL;
