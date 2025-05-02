.PHONY: migrate migrate-checksum help

help:
	@echo "Available commands:"
	@echo "  make migrate NAME=initial_migration  - Create a new migration"
	@echo "  make migrate-checksum               - Regenerate checksum"

migrate:
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Error: NAME parameter is required"; \
		echo "Usage: make migrate NAME=your_migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@go run -mod=mod ./ent/migrate/main.go -name=$(NAME)

migrate-checksum:
	@echo "Regenerating checksum..."
	@go run -mod=mod ./ent/migrate/main.go -checksum