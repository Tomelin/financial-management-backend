
build:
	go build -o server main.go

swag:
	swag init --parseDependency  --parseInternal --parseDepth 1 -g cmd/main.go  --output docs/swagger

tidy:
	go mod tidy

run-swag:
	make swag
	go run cmd/*

run:
	go run cmd/*

sqlc:
	cd infra/database/sqlc \
	&& sqlc generate

protobuf:
	mkdir -p pkg \
	&& cd proto \
	&& buf mod update \
	&& buf generate \
	&& cd ..

mocks_core:
	mockery -r \
	--dir=internal/core/ \
	--output=tests/coremocks \
	--outpkg=coremocks \
	--all

mocks_infra:
	mockery -r \
	--dir=internal/infra/ \
	--output=tests/inframocks \
	--outpkg=inframocks \
	--all

mocks_pkg:
	mockery -r \
	--dir=pkg/ \
	--output=tests/pkgmocks \
	--outpkg=pkgmocks \
	--all
	
mocks:
	make mocks_core
	make mocks_infra
	make mocks_pkg

test:
	go test ./... -cover -v

coverprofile:
	go test ./internal/core/entity/... ./internal/core/service/... -v -coverprofile=/tmp/coverage.out
	go tool cover -func=/tmp/coverage.out
	go tool cover -html=/tmp/coverage.out

cover:
	go test ./internal/core/service -v -cover

watch:
	clear
	ulimit -n 1000 #increase the file watch limit, might required on MacOS
	# make swag
	make tidy
	reflex -s -r '\.go$$' make run
	# watch --chgexit -n 1 "ls --all -l --recursive --full-time | sha256sum"  && echo "Detected the modification of a file or directory"
