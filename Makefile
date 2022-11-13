.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/frk

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf frk frk.exe
