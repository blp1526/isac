.PHONY: all
all: snapshot

.PHONY: vet
vet:
	go vet ./...
	@echo

.PHONY: test
test: vet
	go test -v --cover ./...
	@echo

.PHONY: snapshot
snapshot: test
	goreleaser --rm-dist --snapshot
	@echo
