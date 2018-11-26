# go-twmedia
My practice of Golang

```bash
$ dep ensure
$ go build -ldflags '-w -s' -o bin/twmedia
$ bin/twmedia <URL>

# or
$ go run main.go <URL>

# or
$ docker run --rm -v $PWD:/app mmktomato/go-twmedia <URL>
```
