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

.PHONY: clean
clean:
	rm -rf dist
	@echo

.PHONY: tagging
tagging:
	git tag -a ${TAG} -m "${TAG} release"
	@echo

.PHONY: release
release:
	goreleaser --rm-dist
	@echo
