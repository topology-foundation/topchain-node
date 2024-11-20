if [ -z "../build/topchaind" ]; then
    echo "topchaind not found, building..."
    make build
fi

if [ -z "$(docker images -q topchain-node:latest 2> /dev/null)" ]; then
    echo "Docker image not found, building..."
    docker build -t topchain-node:latest .
fi
