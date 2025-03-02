LDFLAGS=-ldflags "-X main.Version=${VERSION}"

.PHONY: clean cleancaches build updatedeps test

check:
	@test -z $(shell gofmt -l . | tee /dev/stderr) || ( echo "[ERR] Fix formatting issues with 'make fmt'" && exit 1 )
	@test -z $(shell go vet | tee /dev/stderr) || ( echo "[ERR] Fix issues found from 'go vet'" && exit 1 )

build: env-vars
	env GOOS=linux   GOARCH=arm64 go build -C src ${LDFLAGS} -o ../bin/sm-login-linux-arm64
	env GOOS=linux   GOARCH=amd64 go build -C src ${LDFLAGS} -o ../bin/sm-login-linux-amd64
	env GOOS=windows GOARCH=arm64 go build -C src ${LDFLAGS} -o ../bin/sm-login-windows-arm64.exe
	env GOOS=windows GOARCH=amd64 go build -C src ${LDFLAGS} -o ../bin/sm-login-windows-amd64.exe
	env GOOS=darwin  GOARCH=arm64 go build -C src ${LDFLAGS} -o ../bin/sm-login-darwin-arm64
	env GOOS=darwin  GOARCH=amd64 go build -C src ${LDFLAGS} -o ../bin/sm-login-darwin-amd64

clean:
	@rm -rf ./bin

cleancaches:
	@go clean -x -modcache
	@go clean -x -testcache
	@go clean -x -cache

test:
	go test -C src -v -timeout 30s -short -cover ./...


env-vars:
ifndef VERSION
	$(error VERSION is undefined)
endif






