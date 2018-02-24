.PHONY: all
all: build

.PHONY: gometalinter
gometalinter:
	gometalinter --vendor --skip vendor --disable-all \
		--enable vet \
		--enable gofmt \
		--enable golint \
		--enable goimports \
		./...
	@echo

.PHONY: test
test: gometalinter
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
