.PHONY: all
all: build

.PHONY: lint
lint:
	go vet ./...
	gofmt -l .
	@echo

.PHONY: test
test: lint
	go test -v --cover ./...
	@echo

.PHONY: readme
readme: test
	go run _generator/readme.go
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
