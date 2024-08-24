postgres:
	    docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -d postgres

createdb:
		docker exec -it postgres-container createdb --username=postgres --owner=postgres cedar-bank

dropdb:
		docker exec -it postgres-container dropdb --username=postgres cedar-bank -f

migrate:
		migrate -path db/migrations/ -database "postgresql://postgres:password@localhost:5432/cedar-bank?sslmode=disable" -verbose up

sqlc:
		sqlc generate

test:
		go test -v -cover ./...

mock:
		mockgen -destination db/mock/store.go -package mockdb github.com/devphasex/cedar-bank-api/db/sqlc Store


server:
	go run .

.PHONY: createdb dropdb postgres migrate sqlc
