postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postgres-container createdb --username=postgres --owner=postgres cedar-bank
dropdb:
	docker exec -it postgres-container dropdb --username=postgres -f cedar-bank

migrate:
	migrate -path db/migrations/ -database "postgresql://postgres:password@localhost:5432/cedar-bank?sslmode=disable" -verbose up

.PHONY: createdb dropdb postgres migrate
