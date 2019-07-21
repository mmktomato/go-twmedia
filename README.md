# go-twmedia
My practice of Golang


## Usage
```bash
$ make
$ bin/twmedia [URL]

# or
$ go run main.go [URL]

# or
$ docker run --rm -v $PWD:/mnt mmktomato/go-twmedia:latest [URL]
```

## Sptions
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

## Prerequisites
### Runtime
* ffmpeg

### Build
* go 1.12 or higher
* gomock (https://github.com/golang/mock)[https://github.com/golang/mock]
