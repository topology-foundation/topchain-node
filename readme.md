# Mandu Node

Golang implementation of the Mandu node built using Cosmos SDK and Tendermint.

## Local Development Setup

### Pre-requisites

- [Go](https://go.dev/dl/) 1.23 installed
- [Make](https://www.gnu.org/software/make/) installed

### Installation

1. Clone the repository
2. Run `make build` to build the binary at `build/mandud`
3. Run `make config-devnet` to initialize the chain with a mock configuration
4. Run `make run-devnet` to start the node
