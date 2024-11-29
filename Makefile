chain_name?=mandu
home?=$(shell pwd)/build/$(chain_name)

.PHONY: build clean
build:
	go build -o ./build/mandud ./cmd/mandud/main.go

clean:
	rm -rf ./build

config-mock:
	./build/mandud init $(chain_name) --home $(home)
	./build/mandud genesis add-genesis-account alice 100000000stake --home $(home)
	./build/mandud config set app minimum-gas-prices 0mandu --home $(home)
	./build/mandud genesis gentx alice 100000000stake --home $(home)
	./build/mandud genesis collect-gentxs --home $(home)

docker-build:
	docker build -t mandu-node:latest .


## proto-all: Format, lint and generate code from proto files using buf.
proto-all: proto-format proto-lint proto-gen format

## proto-gen: Run buf generate.
proto-gen:
	@echo Generating code from proto...
	@sh scripts/protocgen.sh

## proto-format: Run buf format and update files with invalid proto format>
proto-format:
	@echo Formatting proto files...
	@buf format --write

## proto-lint: Run buf lint.
proto-lint:
	@echo Linting proto files...
	@buf lint

.PHONY: proto-all proto-gen proto-format proto-lint
