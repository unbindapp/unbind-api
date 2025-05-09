-- reverse: create index "github_apps_uuid_key" to table: "github_apps"
DROP INDEX "github_apps_uuid_key";
-- reverse: modify "github_apps" table
ALTER TABLE "github_apps" DROP COLUMN "uuid";
