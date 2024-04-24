postgres:
	docker run --rm -d --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1079 -p 5432:5432 postgres:15

createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres RoommateTap
migarte:
	migrate create -ext sql -dir pkg/database/migration -seq init_schema

migrateup:
	migrate -path pkg/database/migration -database "postgresql://postgres:1079@localhost:5432/RoommateTap?sslmode=disable" -verbose up

migratedown:
	migrate -path  pkg/database/migration -database "postgresql://postgres:1079@localhost:5432/RoommateTap?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres15 dropdb -U postgres RoommateTap
#
#proto:
#	protoc ./internal/proto/*.proto --go_out=. --go-grpc_out=.

redis:
	docker run -d -p 6379:6379 --name my-redis-container redis


.PHONY: postgres createdb dropdb