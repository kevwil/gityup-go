.PHONY: help all vendor build test run

APP_NAME=gityup-$(shell go env GOOS)-$(shell go env GOARCH)

## all           : run default tasks [vendor lint test]
all: vendor lint test

## clean         : clean
clean:
	@echo "\n*** go clean stuff\n"
	@go clean -i -v .
	@rm -f gityup-*

## vendors       : go mod dependencies
vendor: go.mod
	@echo "\n*** update go vendor folder and modules\n"
	@go mod vendor
	@go mod verify
	@go mod tidy -v -compat=1.25

## lint          : syntax checking and formatting
lint:
	@echo "\n*** lint and fmt the code\n"
	@goimports -w .
	@staticcheck ./...
	@golangci-lint run --no-config -E goconst -E prealloc -E unparam
	@gosec -quiet ./...
	@golangci-lint fmt --no-config -E gofmt

## build         : build production container image
build: go.mod
	@echo "\n*** build app\n"
	@go build -mod vendor -a -o $(APP_NAME) .

## test          : run tests
test: go.mod
	@echo "\n*** run tests\n"
	@go test -mod vendor -v .

_vuln_code: go.mod
	@echo "\n*** vuln check on code base\n"
	@govulncheck

_vuln_binary: $(APP_NAME)
	@echo "\n*** vuln check on binary file\n"
	@govulncheck -mode binary -show verbose $(APP_NAME)

## vuln          : run govulncheck on code and binary
vuln: _vuln_code _vuln_binary

## run           : run the app
run: go.mod
	@echo "\n*** run the app using local code (not a binary)\n"
	@go run .

## help          : print self-documented target info from this file
help : Makefile
	@sed -n 's/^##//p' $<
