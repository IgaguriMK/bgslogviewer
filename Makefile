GOLINT_OPTS=-min_confidence 1.0 -set_exit_status

VERSION=0.1.0
IMAGENAME=bgslogviewer:$(VERSION)

.PHONY: all
all: dep build lint

.PHONY: b
b: build lint

.PHONY: build
build:
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
dr: docker-build docker-run

.PHONY: docker-build
docker-build: build-linux
	docker build -t $(IMAGENAME) .

.PHONY: docker-run
docker-run:
	docker run -it --rm -p 8080:8080 $(IMAGENAME)


.PHONY: clean
clean:
	- rm bgslogviewer
	- rm bgslogviewer.exe
