.PHONY: lint
lint:
	golangci-lint run

.PHONY: slave
slave:
	./scripts/main slave

.PHONY: master
master:
	./scripts/main master

.PHONY: generate
generate:
	./scripts/gen_proto
	go generate ./...

.PHONY: test
test:
	go test ./... -cover