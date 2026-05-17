DB_URL=postgres://admin:secret@localhost:5433/chatdb?sslmode=disable
MIGRATIONS_DIR=./db/migrations

# ======================
# MIGRATIONS (GOOSE)
# ======================

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status


# ======================
# POSTGRES CONTAINER
# ======================

postgresinit:
	docker run -d --name postgredb \
	-e POSTGRES_USER=admin \
	-e POSTGRES_PASSWORD=secret \
	-p 5433:5432 \
	-v pgdata:/var/lib/postgresql/data \
	postgres:16-alpine

postgres:
	docker exec -it postgredb psql -U admin -d chatdb

createdb:
	docker exec -it postgredb createdb -U admin chatdb

dropdb:
	docker exec -it postgredb dropdb -U admin --if-exists chatdb

logs:
	docker logs -f postgredb

stopdb:
	docker stop postgredb

startdb:
	docker start postgredb

removedb:
	docker rm -f postgredb


# ======================
# PHONY
# ======================

.PHONY: \
migrate-up migrate-down migrate-reset migrate-create migrate-status \
postgresinit postgres createdb dropdb logs stopdb startdb removedb