GOLINT_OPTS := -min_confidence 0.6 -set_exit_status

VERSION   := $(shell cat .version)
IMAGENAME := bgslogviewer:$(VERSION)

.PHONY: all
all: clean dep build lint docker-build

.PHONY: b
b: build lint

.PHONY: build
build:
	echo $(IMAGENAME)
	go build bgslogviewer.go

.PHONY: build-linux
build-linux:
	env GOOS=linux GOARCH=amd64 go build -o bgslogviewer bgslogviewer.go

.PHONY: lint
lint:
	golint $(GOLINT_OPTS) bgslogviewer.go

.PHONY: dep
dep:
	dep ensure


.PHONY: dr
dr: clean docker-build docker-run

.PHONY: docker-build
docker-build: build-linux lint
	docker build --no-cache -t $(IMAGENAME) .

.PHONY: docker-run
docker-run:
	docker-compose up


.PHONY: clean
clean:
	- rm bgslogviewer
	- rm bgslogviewer.exe
	- rm log/*
	- docker-compose down
