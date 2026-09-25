default: lint

lint:
	golangci-lint run --config .golangci.yml -v ./pkg/errformatter/

.PHONY: lint