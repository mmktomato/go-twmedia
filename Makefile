build:
	dep ensure -v
	go build -ldflags '-w -s' -o bin/twmedia

goinstall:
	dep ensure -v
	go install -ldflags='-w -s'

clean:
	rm -rf bin/*
	rm -rf vendor/*
