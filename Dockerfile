FROM golang:1.23 as build

WORKDIR /topchain-build
COPY . .
RUN go mod tidy
RUN go build -o topchaind ./cmd/topchaind

FROM ubuntu:latest as runtime
COPY --from=build /topchain-build/topchaind /usr/local/bin/topchaind

ENTRYPOINT ["topchaind"]
