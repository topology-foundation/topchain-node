build:
	go build -o ./build/topchaind ./cmd/topchaind/main.go

container:
	docker build -t topchain-node:latest .

chain_name?=topchain
home?=$(shell pwd)/build/$(chain_name)
config-mock:
	./build/topchaind init $(chain_name) --home $(home)
	./build/topchaind genesis add-genesis-account alice 100000000stake --home $(home)
	./build/topchaind config set app minimum-gas-prices 0top --home $(home)
	./build/topchaind genesis gentx alice 100000000stake --home $(home)
	./build/topchaind genesis collect-gentxs --home $(home)
