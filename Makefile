BINARY = stripe-copy

SRCS = %.go

all: $(BINARY)

$(BINARY): **/**/*.go **/*.go *.go
	go build -ldflags "-X main.buildCommit `git describe --long --tags --dirty --always`" .

deps:
	go get .

build: $(BINARY)

clean:
	rm $(BINARY)

test:
	go test ./...
	golint ./...
