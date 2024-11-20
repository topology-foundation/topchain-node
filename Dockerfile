FROM golang:1.23

COPY ./build/topchaind /usr/local/bin/topchaind
ENTRYPOINT ["topchaind"]
