-- create "service_groups" table
CREATE TABLE "service_groups" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "environment_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "service_groups_environments_service_groups" FOREIGN KEY ("environment_id") REFERENCES "environments" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- modify "services" table
ALTER TABLE "services" ADD COLUMN "service_group_id" uuid NULL, ADD
CONSTRAINT "services_service_groups_services" FOREIGN KEY ("service_group_id") REFERENCES "service_groups" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
