-- modify "deployments" table
ALTER TABLE "deployments" ADD COLUMN "queued_at" timestamptz NULL;
