postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb --username=root --owner=root task_tracker 

dropdb:
	docker exec -it postgres17 dropdb --username=root task_tracker 

migrateup:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/task_tracker?sslmode=disable" -verbose up

migratedown:
	migrate -path internal/db/migration -database "postgresql://root:secret@localhost:5432/task_tracker?sslmode=disable" -verbose down 

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown test server