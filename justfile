default:
    @just --list

build:
    go build -v

test *ARGS:
    ginkgo -r {{ARGS}}

format:
    gofmt -s -w .

lint:
    go vet ./...

deps:
    go get

deps-update:
    go get -u ./...
    go mod tidy

deps-outdated:
    go list -u -m all

clean:
    rm --force dfm

install-tools:
    go install github.com/onsi/ginkgo/ginkgo@v1.16.5
