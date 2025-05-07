-- reverse: modify "services" table
ALTER TABLE "services" DROP CONSTRAINT "services_templates_services", DROP COLUMN "template_id";
-- reverse: modify "templates" table
ALTER TABLE "templates" DROP COLUMN "immutable";
