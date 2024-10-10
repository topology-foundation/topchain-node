FROM golang:1.23

WORKDIR /topchain-node
COPY . .

RUN go mod tidy
RUN go build cmd/topchaind/main.go && \
	mv main /usr/local/bin/topchaind

RUN chmod +x /topchain-node/entrypoint.sh

ENTRYPOINT ["/bin/bash", "/topchain-node/entrypoint.sh"]
