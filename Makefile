GOLINT_OPTS=-min_confidence 1.0 -set_exit_status


.PHONY: all
all: dep build lint


.PHONY: build
build:
	go build bgslogviewer.go

.PHONY: lint
lint:
	golint $(GOLINT_OPTS) bgslogviewer.go


.PHONY: dep
dep:
	dep ensure


