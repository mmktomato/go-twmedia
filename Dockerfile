FROM golang:1.11.2-alpine3.8 as builder

# ENV GOPATH /go
WORKDIR $GOPATH/src/github.com/mmktomato/go-twmedia

COPY . .

RUN apk update \
    && apk add --no-cache git dep make \
    && make goinstall \
    && make clean \
    && apk del --purge git dep make

###################

FROM alpine:3.8

WORKDIR /mnt

COPY --from=builder /go/bin/go-twmedia /usr/local/bin/

RUN apk update && apk add --no-cache ffmpeg

ENTRYPOINT ["go-twmedia"]
