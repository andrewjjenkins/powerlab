GOCMD=go
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOCLEAN=$(GOCMD) clean

all: build

deps:
	#go get ./...
	#go install "github.com/go-swagger/go-swagger/cmd/swagger@latest"
	cd client && npm install

generate:
	$(GOGENERATE) ./cmd/... ./pkg/...
	cd client && npm run generate	

run: deps
	cd client && npm run dev

build: generate
	$(GOBUILD) ./pkg/...
	$(GOBUILD) ./cmd/powerlab
	cd client && npm run build

clean:
	rm -rf client/src/generated/* pkg/serve/api powerlab
	$(GOCLEAN) ./pkg/... ./cmd/...

.PHONY: all deps generate build load clean
