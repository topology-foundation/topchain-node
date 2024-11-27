FROM golang:1.23 as build

WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o mandud ./cmd/mandud

FROM ubuntu:latest as runner
COPY --from=build /build/mandud /usr/local/bin/mandud

ENTRYPOINT ["mandud"]
