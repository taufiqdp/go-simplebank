migrateup:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: migrateup migrateup