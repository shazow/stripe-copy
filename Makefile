BINARY = stripe-copy

all: $(BINARY)

$(BINARY): *.go
	go build -ldflags "-X main.version $(git describe --long --tags --dirty --always)" .

deps:
	go get .

build: $(BINARY)

run: $(BINARY)
	./$(BINARY) -vv --pretend

clean:
	rm $(BINARY)

test:
	go test ./...
	golint ./...
