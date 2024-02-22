sqli:
	sqlc init

sqlg:		
	sqlc generate


server:
	go run ./cmd/api

test:
	go test -v -count=1 ./...

.PHONY: sqli sqlg server test 
