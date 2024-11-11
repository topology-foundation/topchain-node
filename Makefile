build-mac-amd64:
	GOOS=darwin GOARCH=amd64 go build -o ./build/topchaind ./cmd/topchaind/main.go

build-mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o ./build/topchaind ./cmd/topchaind/main.go

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o ./build/topchaind ./cmd/topchaind/main.go

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o ./build/topchaind ./cmd/topchaind/main.go

build-all: build-mac-amd64 build-mac-arm64 build-linux-amd64 build-linux-arm64
