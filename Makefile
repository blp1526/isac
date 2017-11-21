.PHONY: all
all: test

.PHONY: vet
vet:
	go vet ./...
	@echo

.PHONY: test
test: vet
	go test -v --cover ./...
	@echo
