-- +goose Up
-- modify "service_configs" table
ALTER TABLE "service_configs" 
  RENAME COLUMN "docker_builder_path" TO "docker_builder_dockerfile_path",
  RENAME COLUMN "docker_builder_context" TO "docker_builder_build_context",
  RENAME COLUMN "install_command" TO "railpack_builder_install_command",
  RENAME COLUMN "build_command" TO "railpack_builder_build_command";

-- +goose Down
-- reverse: modify "service_configs" table
ALTER TABLE "service_configs" 
  RENAME COLUMN "railpack_builder_build_command" TO "build_command",
  RENAME COLUMN "railpack_builder_install_command" TO "install_command",
  RENAME COLUMN "docker_builder_build_context" TO "docker_builder_context",
  RENAME COLUMN "docker_builder_dockerfile_path" TO "docker_builder_path";