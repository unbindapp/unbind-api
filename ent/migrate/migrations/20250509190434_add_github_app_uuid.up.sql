-- modify "github_apps" table
ALTER TABLE "github_apps" ADD COLUMN "uuid" uuid NOT NULL;
-- create index "github_apps_uuid_key" to table: "github_apps"
CREATE UNIQUE INDEX "github_apps_uuid_key" ON "github_apps" ("uuid");
