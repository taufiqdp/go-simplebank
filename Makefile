migrateup:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgres://root:pwd@127.0.0.1:5432/simplebank?sslmode=disable" -verbose down 1
 
test:
	go clean -testcache && go test -v -cover ./...

sqlc:
	sqlc generate

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go github.com/taufiqdp/go-simplebank/db/sqlc Store

.PHONY: migrateup migrateup test sqlc server mock