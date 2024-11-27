#!/bin/bash
./scripts/clean.sh

if [ -z "$(docker images -q mandu-node:latest 2> /dev/null)" ]; then
    echo "Docker image not found, please run make container first"
    exit 0
fi


if [ -z "$HOME" ]; then
    echo "HOME is not set, defaulting to ."
    HOME="."
fi

genesis_dirs=("val-alice" "val-bob" "node-carol")
validator_dirs=("val-alice" "val-bob")
for dir in "${genesis_dirs[@]}"; do
    echo "Creating directory: $dir"
    mkdir -p "$HOME/mandu-devnet/$dir"
done

# provide:
# - directory
# - command to run, space separated
run_command_in_docker() {
    docker run -i --rm -v "$HOME/mandu-devnet/$1:/root/.mandu" mandu-node:latest ${@:2}
}

for dir in "${genesis_dirs[@]}"; do
    run_command_in_docker $dir init mandu --chain-id mandu
done

for dir in "${validator_dirs[@]}"; do
    name=$(echo $dir | cut -d'-' -f2)
    run_command_in_docker $dir keys --keyring-backend file --keyring-dir /root/.mandu/keys add $name <<< "12345678"$'\n'"12345678"
done

ALICE=$(echo "12345678" | run_command_in_docker val-alice keys --keyring-backend file --keyring-dir /root/.mandu/keys show alice --address)
BOB=$(echo "12345678" | run_command_in_docker val-bob keys --keyring-backend file --keyring-dir /root/.mandu/keys show bob --address)
echo "ALICE: $ALICE"
echo "BOB: $BOB"

# https://tutorials.cosmos.network/hands-on-exercise/4-run-in-prod/1-run-prod-docker.html#genesis
run_command_in_docker val-alice genesis add-genesis-account $ALICE 100000000000stake,10000000000utop
mv "$HOME/mandu-devnet/val-alice/config/genesis.json" "$HOME/mandu-devnet/val-bob/config/genesis.json"
run_command_in_docker val-bob genesis add-genesis-account $BOB 100000000000stake,10000000000utop
echo "12345678" | run_command_in_docker val-bob genesis gentx bob 1000000stake --chain-id mandu --keyring-backend file --keyring-dir /root/.mandu/keys --account-number 0 --sequence 0 --gas 100000 --gas-prices 0.025stake
mv "$HOME/mandu-devnet/val-bob/config/genesis.json" "$HOME/mandu-devnet/val-alice/config/genesis.json"
echo "12345678" | run_command_in_docker val-alice genesis gentx alice 1000000stake --chain-id mandu --keyring-backend file --keyring-dir /root/.mandu/keys --account-number 0 --sequence 0 --gas 100000 --gas-prices 0.025stake
cp "$HOME/mandu-devnet/val-bob/config/gentx/"* "$HOME/mandu-devnet/val-alice/config/gentx"
run_command_in_docker val-alice genesis collect-gentxs
run_command_in_docker val-alice genesis validate-genesis

# distribute genesis.json
for dir in "${genesis_dirs[@]}"; do
    cp "$HOME/mandu-devnet/val-alice/config/genesis.json" "$HOME/mandu-devnet/$dir/config/genesis.json"
done

for dir in "${genesis_dirs[@]}"; do
    seeds=""
    for other_dir in "${genesis_dirs[@]}"; do
        if [ "$dir" != "$other_dir" ]; then
            address=$(run_command_in_docker "$other_dir" tendermint show-node-id)@$other_dir:26656
            if [ -n "$seeds" ]; then
                seeds="$seeds,$address"
            else
                seeds="$address"
            fi
        fi
    done

    # Update the seeds field in config.toml
    docker run --rm -i -v "$HOME/mandu-devnet/$dir:/root/.mandu" --entrypoint sed mandu-node:latest \
        -Ei "s/(seeds *= *\"[^\"]*)\"/\1$seeds\"/" "/root/.mandu/config/config.toml"

    run_command_in_docker $dir config set app minimum-gas-prices 0.025stake
done

# RPC on all addresses
docker run --rm -i -v "$HOME/mandu-devnet/node-carol:/root/.mandu" --entrypoint sed mandu-node:latest -Ei '0,/^laddr = .*$/{s/^laddr = .*$/laddr = "tcp:\/\/0.0.0.0:26657"/}' /root/.mandu/config/config.toml

# start nodes
docker-compose -f ./scripts/docker-compose.yml up -d
