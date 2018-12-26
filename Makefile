.PHONY: build
build:
	dep ensure -v
	go build -ldflags '-w -s' -o bin/twmedia

.PHONY: goinstall
goinstall:
	dep ensure -v
	go install -ldflags='-w -s'

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: mockgen
mockgen:
	mockgen -source svc/extcmd/extcmd.go -destination svc/extcmd/_mock/mock_extcmd.go

