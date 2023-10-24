postgres:
	docker run --name postgresdb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
stopdbcontainer:
	docker stop postgresdb
startdbcontainer:
	docker start postgresdb
createdb: 
	docker exec -it postgresdb createdb --username=root --owner=root production
dropdb:
	docker exec -it postgresdb dropdb production
migrateup:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5432/production?sslmode=disable" --verbose up
migratedown:
	migrate --path db/migration --database "postgresql://root:secret@localhost:5432/production?sslmode=disable" --verbose down
sqlc:
	sqlc generate
server:
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown sqlc stopdbcontainer startdbcontainer server