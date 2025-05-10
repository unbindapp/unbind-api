-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" DROP COLUMN "build_command", DROP COLUMN "install_command";
