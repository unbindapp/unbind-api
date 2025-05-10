-- modify "service_configs" table
ALTER TABLE "service_configs" ADD COLUMN "install_command" character varying NULL, ADD COLUMN "build_command" character varying NULL;
