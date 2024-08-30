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
		go test -v -cover  -count 1 ./...

mock:
		mockgen -destination db/mock/store.go -package mockdb github.com/devphasex/cedar-bank-api/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
		proto/*.proto
	statik -src=./doc/swagger -dest=./doc

server:
	go run .

.PHONY: createdb dropdb postgres migrate sqlc proto
