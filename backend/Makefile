# ======================
# CONFIG
# ======================

DB_URL=postgres://admin:secret@localhost:5433/chatdb?sslmode=disable
MIGRATIONS_DIR=./db/migrations


# ======================
# DEVELOPMENT (AIR HOT RELOAD)
# ======================

dev:
	docker compose up app-dev

dev-build:
	docker compose up -d --build app-dev

dev-down:
	docker compose stop app-dev

dev-logs:
	docker compose logs -f app-dev

# ======================
# DOCKER COMPOSE CORE
# ======================

up:
	docker compose up -d --build

down:
	docker compose down

restart:
	docker compose down && docker compose up -d --build

logs:
	docker compose logs -f

ps:
	docker compose ps

build:
	docker compose build

# ======================
# APP CONTROL
# ======================

app:
	docker compose up -d app

app-logs:
	docker compose logs -f app

# ======================
# POSTGRES
# ======================

db:
	docker compose exec postgres psql -U admin -d chatdb

db-shell:
	docker compose exec postgres sh

db-reset:
	docker compose exec postgres dropdb -U admin --if-exists chatdb
	docker compose exec postgres createdb -U admin chatdb

db-logs:
	docker compose logs -f postgres

# ======================
# REDIS
# ======================

redis-cli:
	docker compose exec redis redis-cli -a redis090517

redis-logs:
	docker compose logs -f redis

# ======================
# MIGRATIONS (GOOSE via migrate container)
# ======================

migrate-up:
	docker compose run --rm migrate goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	docker compose run --rm migrate goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-reset:
	docker compose run --rm migrate goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

migrate-status:
	docker compose run --rm migrate goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-create:
	docker compose run --rm migrate goose -dir $(MIGRATIONS_DIR) create $(name) sql

# ======================
# DATA MANAGEMENT
# ======================

db-clean-data:
	@echo "Wiping data from all tables (preserving schema)..."
	docker cp ./db/scripts/clean_data.sql relio-postgres:/tmp/clean.sql
	docker exec -i relio-postgres psql -U admin -d chatdb -f /tmp/clean.sql
	docker exec -i relio-postgres rm -f /tmp/clean.sql

# ======================
# UTILITIES
# ======================

status:
	docker compose ps

stop:
	docker compose stop

start:
	docker compose start

# ======================
# CLEANUP
# ======================

clean:
	docker compose down -v
	docker system prune -f

# ======================
# PHONY
# ======================

.PHONY: \
up down restart logs ps build \
app app-logs \
db db-shell db-reset db-logs \
redis-cli redis-logs \
migrate-up migrate-down migrate-reset migrate-status migrate-create \
db-clean-data \
status stop start clean