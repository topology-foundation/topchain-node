FROM golang:1.23

WORKDIR /topchain-node
COPY . .

RUN curl https://get.ignite.com/cli | bash && \
	chmod +x ignite

RUN ./ignite chain build && \
	./ignite chain init
RUN chmod +x /topchain-node/entrypoint.sh

ENTRYPOINT ["/bin/bash", "/topchain-node/entrypoint.sh"]
