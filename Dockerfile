# Use an official golang as a parent image
FROM golang:1.23

WORKDIR /topchain-node

COPY . .

RUN go mod tidy

# Install Ignite CLI
RUN curl https://get.ignite.com/cli | bash && \
    mv ignite /usr/local/bin/ignite && \
    chmod +x /usr/local/bin/ignite

# Command to serve the Topchain node
CMD ["/bin/sh", "-c", "echo 'No' | ignite chain serve"]
