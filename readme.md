# TopChain Installation Guide

This file provides instructions for installing TopChain both manually and using Docker.

## Manual Installation

### 1. Install Ignite CLI

Ignite CLI is a powerful tool for Cosmos SDK-based blockchain development. To install it:

a. Install Ignite CLI:
   ```
   curl https://get.ignite.com/cli! | bash
   ```
   If you are having troubles installing ignite, please refer to the [official documentation](https://docs.ignite.com/cli/install) for more details.

b. Verify the installation:
   ```
   ignite version
   ```

### 2. Clone the TopChain repository:
   ```
   git clone https://github.com/topology-foundation/topchain-node.git
   cd topchain-node
   ```

### 3. Build the TopChain binary:
   ```
   ignite chain build
   ```

### 4. Verify the installation:
   ```
   topchaind version
   ```
   Note: Ignite installs the binary in the `~/go/bin/` directory. Ensure this directory is included in your PATH.

### 5. Initialize the TopChain node:
   ```
   topchaind init <your_node_name> --chain-id <topchain_chain_id>
   ```

### 6. Start the TopChain node:
   ```
   topchaind start
   ```

## Docker-based Installation

1. Clone the TopChain repository:
   ```
   git clone https://github.com/topology-foundation/topchain-node.git
   cd topchain-node
   ```

2. Build the Docker image:
   ```
   docker build -t topchain .
   ```

3. Run the TopChain node in a Docker container:
   ```
   docker run -d --name topchain-node topchain
   ```

4. Check the logs of the running container:
   ```
   docker logs -f topchain-node
   ```

