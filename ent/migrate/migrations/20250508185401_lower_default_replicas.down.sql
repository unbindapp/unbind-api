-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" ALTER COLUMN "replicas" SET DEFAULT 2;
