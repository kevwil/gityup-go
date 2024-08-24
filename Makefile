.PHONY: help all vendor build test run

APP_NAME=gityup-$(shell go env GOOS)-$(shell go env GOARCH)

## all           : run default tasks [clean vendor lint test]
all: clean vendor lint test

## clean         : clean
clean:
	@echo "\n*** go clean stuff\n"
	@go clean -i -v .
	@rm -f $(APP_NAME)

## vendors       : go mod dependencies
vendor: go.mod
	@echo "\n*** update go vendor folder and modules\n"
	@go mod vendor
	@go mod verify
	@go mod tidy -v -compat=1.23

## lint          : syntax checking and formatting
lint:
	@echo "\n*** lint and fmt the code\n"
	@goimports -w .
	@golangci-lint run --no-config -D gosimple -D staticcheck -D unused -D govet -D typecheck -E goconst -E gofmt -E goimports -E gosec -E prealloc -E unparam

## build         : build production container image
build: go.mod
	@echo "\n*** build app\n"
	@go build -mod vendor -a -o $(APP_NAME) .

## test          : run tests
test:
	@echo "\n*** run tests\n"
	@go test -mod vendor -v .

## run           : run the app
run:
	@echo "\n*** run the app\n"
	@go run .

## help          : print self-documented target info from this file
help : Makefile
	@sed -n 's/^##//p' $<
