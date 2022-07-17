POSTGRESQL_URL := 'pgx://postgres:qazwsxedc@10.10.152.12:5432/db_test?sslmode=disable&search_path=public'
PATH_PROTO := ./api/calendar/server.proto

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

gen-grpc:
	protoc -I . --go_out . --go-grpc_out . $(PATH_PROTO)

gen-gateway:
	protoc -I . \
	--grpc-gateway_out . \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt generate_unbound_methods=true \
	$(PATH_PROTO)

gen-swagger:
	protoc -I . --swagger_out ./internal/server/grpc $(PATH_PROTO)

gen-all:
	protoc -I . \
	--go_out $(PATH_GEN_GOLANG_GRPC) \
	--go_opt paths=source_relative \
	--go-grpc_out $(PATH_GEN_GOLANG_GRPC) \
	--go-grpc_opt paths=source_relative \
	--grpc-gateway_out $(PATH_GEN_GOLANG_GRPC) \
	--grpc-gateway_opt logtostderr=true \
	--grpc-gateway_opt paths=source_relative \
	--grpc-gateway_opt generate_unbound_methods=true \
	--swagger_out $(PATH_GEN_GOLANG_GRPC) \
	./api/calendar/server.proto

generate:
	protoc -I . --go_out . --go-grpc_out . $(PATH_PROTO)

migrate-up:
	migrate -database $(POSTGRESQL_URL) -path ./migrate/calendar up

migrate-down:
	migrate -database $(POSTGRESQL_URL) -path ./migrate/calendar down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.46.2	

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
