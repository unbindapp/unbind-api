-- modify "templates" table
ALTER TABLE "templates" ADD COLUMN "immutable" boolean NOT NULL DEFAULT false;
-- modify "services" table
ALTER TABLE "services" ADD COLUMN "template_id" uuid NULL, ADD
CONSTRAINT "services_templates_services" FOREIGN KEY ("template_id") REFERENCES "templates" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
