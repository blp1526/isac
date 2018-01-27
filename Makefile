.PHONY: all
all: build

.PHONY: vet
vet:
	go vet ./...
	@echo

.PHONY: test
test: vet
	go test -v --cover ./...
	@echo

.PHONY: readme
readme: vet
	go run _doc/readme.go
	@echo

.PHONY: build
build: readme
	rm -rf dist/
	mkdir -p dist/
	go build -o dist/isac
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
