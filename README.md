# go-twmedia
My practice of Golang

```bash
$ make
$ bin/twmedia [URL]

# or
$ go run main.go [URL]

# or
$ docker run --rm -v $PWD:/mnt mmktomato/go-twmedia:latest [URL]
```

## options
### -H, --header
Same as curl's one.

```bash
$ go run main.go [URL] -H "Referer: http://localhost"
```

### -v, --verbose
Shows verbose log.

```bash
$ go run main.go [URL] -v
```
