# go-twmedia
My practice of Golang

```bash
$ dep ensure
$ go build -o bin/twmedia
$ bin/twmedia <URL>

# or
$ go run main.go <URL>
```

```bash
$ docker build -t my/twmedia .
$ docker run --rm -v $PWD:/app <URL>
