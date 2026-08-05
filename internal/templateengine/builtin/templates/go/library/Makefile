.PHONY: check fmt test vet

check: fmt test vet

fmt:
	test -z "$$(gofmt -l .)"

test:
	go test -race ./...

vet:
	go vet ./...
