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
