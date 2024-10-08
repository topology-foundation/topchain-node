# Use an official golang as a parent image
FROM golang:1.23

# Install Ignite CLI
RUN curl https://get.ignite.com/cli | bash && \
    mv ignite /usr/local/bin/ignite && \
    chmod +x /usr/local/bin/ignite

# Clone the Topchain node repository
RUN git clone https://github.com/topology-foundation/topchain-node.git /topchain-node

# Set working directory
WORKDIR /topchain-node

RUN go mod tidy

# Command to serve the Topchain node
CMD ["/bin/sh", "-c", "echo 'No' | ignite chain serve"]
