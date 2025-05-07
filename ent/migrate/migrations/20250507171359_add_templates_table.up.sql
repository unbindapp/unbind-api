-- create "templates" table
CREATE TABLE "templates" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "name" character varying NOT NULL,
  "version" bigint NOT NULL,
  "definition" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "template_name_version" to table: "templates"
CREATE UNIQUE INDEX "template_name_version" ON "templates" ("name", "version");
