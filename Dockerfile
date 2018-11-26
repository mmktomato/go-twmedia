FROM golang:1.11.2-alpine3.8

# ENV GOPATH /go
WORKDIR $GOPATH/src/github.com/mmktomato/go-twmedia

COPY . .

RUN apk update \
    && apk add --no-cache ffmpeg git dep \
    #
    && dep ensure -v \
    && go install -ldflags='-w -s' \
    #
    && apk del --purge git dep

WORKDIR /app
ENTRYPOINT ["go-twmedia"]
