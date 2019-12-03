.PHONY: help run build

all:
	@echo "****************************"
	@echo "**       build tool       **"
	@echo "****************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run API in dev mode"
	@echo "  build                 - build api into bin/ directory"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)

run:
	@(go run cmd/sustainable/main.go)

build:
	@mkdir -p ./bin
	go build -gcflags=-trimpath=${GOPATH} -asmflags=-trimpath=${GOPATH} -i -o ./bin/sustainable ./cmd/sustainable/main.go
