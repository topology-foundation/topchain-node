# topchain node

Golang implementation of the topchain node built using Cosmos SDK

## Local Development Setup
### Pre-requisites
- [Go](https://go.dev/dl/) 1.23 installed
- [Make](https://www.gnu.org/software/make/) installed

### Installation
1. Clone the repository
2. Run `make build` to build the binary at `build/topchaind`
3. Run `make chain_name=your_chain_name home=your_home_dir config-mock` to initialize the chain with a mock configuration
4. Run `./build/topchaind start --home=your_home_dir` to start the node
