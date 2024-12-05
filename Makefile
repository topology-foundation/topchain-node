bin?=mandud
chain_name?=mandu
home?=$(shell pwd)/build/$(chain_name)

.PHONY: build clean
build:
	go build -o ./build/$(bin) ./cmd/mandud/main.go

clean:
	rm -rf ./build

config-devnet:
	rm -rf $(home)
	./build/$(bin) init $(chain_name) --home $(home)
	./build/$(bin) keys add alice --home $(home)
	./build/$(bin) genesis add-genesis-account alice 100000000stake --home $(home)
	./build/$(bin) config set app minimum-gas-prices 0mandu --home $(home)
	./build/$(bin) genesis gentx alice 100000000stake --home $(home)
	./build/$(bin) genesis collect-gentxs --home $(home)

run-devnet:
	./build/$(bin) start --home $(home)

docker-build:
	docker build -t mandu-node:latest .

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo Generating code from proto...
	@sh scripts/protocgen.sh

proto-format:
	@echo Formatting proto files...
	@buf format --write

proto-lint:
	@echo Linting proto files...
	@buf lint

.PHONY: proto-all proto-gen proto-format proto-lint
