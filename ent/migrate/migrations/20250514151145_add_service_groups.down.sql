-- reverse: modify "services" table
ALTER TABLE "services" DROP CONSTRAINT "services_service_groups_services", DROP COLUMN "service_group_id";
-- reverse: create "service_groups" table
DROP TABLE "service_groups";
