FROM golang:1.12-alpine3.10 as builder

# ENV GOPATH /go
WORKDIR $GOPATH/src/github.com/mmktomato/go-twmedia

COPY . .

RUN apk update \
    && apk add --no-cache git dep make \
    && make goinstall \
    && make clean \
    && apk del --purge git dep make

###################

FROM alpine:3.10

WORKDIR /mnt

COPY --from=builder /go/bin/go-twmedia /usr/local/bin/

RUN apk update && apk add --no-cache ffmpeg ca-certificates

ENTRYPOINT ["go-twmedia"]
