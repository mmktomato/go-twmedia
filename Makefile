.PHONY: build
build:
	GO111MODULE=on go build -ldflags '-w -s' -o bin/twmedia

.PHONY: goinstall
goinstall:
	GO111MODULE=on go install -ldflags='-w -s'

.PHONY: clean
clean:
	rm -rf bin/*

.PHONY: mockgen
mockgen:
	mockgen -source svc/extcmd/extcmd.go -destination svc/extcmd/_mock/mock_extcmd.go

