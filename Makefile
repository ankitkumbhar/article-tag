# atomic: Concurrent access to the same coverage counters is guaranteed to be 
# executed one at a time, avoiding race conditions.
COVER_MODE=atomic # other options [count, set]
COVER_PROFILE=coverage.txt # coverage profile, write out file
COVER_REPORT=coverage.html # view the coverage report(.html file) in the browser

GO_TEST_PKG=$(shell go list ./... | grep -v docs | grep -v cmd | grep -v internal/mocks | grep -v internal/response) # grep to ignore
GO_COVER_PKG=$(shell go list ./... | grep -v docs | grep -v cmd | grep -v internal/mocks | grep -v internal/response | paste -sd "," -)

## Compile and execute code
local-run:
	go run cmd/main.go

run:
	docker-compose up -d

stop:
	docker-compose stop

mod:
	go mod tidy && go mod vendor

mocks:
	rm -rf ./internal/mocks
	mockery --all --case underscore --with-expecter --exported --srcpkg ./internal/model --output ./internal/mocks

unit-test:
	go test ${GO_TEST_PKG} -mod=readonly -cover -covermode=${COVER_MODE} -coverprofile=${COVER_PROFILE} -coverpkg=${GO_COVER_PKG}

cover:
	go tool cover -func coverage.txt