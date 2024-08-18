GOCMD=go
GOBUILD=$(GOCMD) build
GOGENERATE=$(GOCMD) generate
GOCLEAN=$(GOCMD) clean

all: build

deps:
	cd pkg/serve/client && npm install

generate:
	$(GOGENERATE) ./cmd/... ./pkg/...
	cd pkg/serve/client && npm run generate	

run: deps
	cd pkg/serve/client && npm run dev

build: generate
	cd pkg/serve/client && npm run build
	$(GOBUILD) ./pkg/...
	$(GOBUILD) ./cmd/powerlab

clean:
	rm -rf pkg/serve/client/src/generated/* pkg/serve/api powerlab
	$(GOCLEAN) ./pkg/... ./cmd/...

.PHONY: all deps generate build load clean
