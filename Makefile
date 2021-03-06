.PHONY: all test clean man glide fast release install version build

GO15VENDOREXPERIMENT=1

PROG_NAME := "cxdig"

all: deps test build install version

build: deps
	@go build -ldflags "-X github.com/sniperkit/cxdig/pkg/cmd.softwareVersion=`cat VERSION`" -o ./bin/$(PROG_NAME) ./cmd/$(PROG_NAME)/*.go

version:
	@which $(PROG_NAME)
	@$(PROG_NAME) version

install: deps
	@go install -ldflags "-X github.com/sniperkit/cxdig/pkg/cmd.softwareVersion=`cat VERSION`" ./cmd/$(PROG_NAME)
	@$(PROG_NAME) version

fast: deps
	@go build -i -ldflags "-X github.com/sniperkit/cxdig/pkg/cmd.softwareVersion=`cat VERSION`-dev" -o ./bin/$(PROG_NAME) ./cmd/$(PROG_NAME)/*.go
	@$(PROG_NAME) version

deps:
	@glide install --strip-vendor

test:
	@go test ./pkg/...

clean:
	@go clean
	@rm -fr ./bin
	@rm -fr ./dist

release: $(PROG_NAME)
	@git tag -a `cat VERSION`
	@git push origin `cat VERSION`

cover:
	@rm -rf *.coverprofile
	@go test -coverprofile=$(PROG_NAME).coverprofile ./pkg/...
	@gover
	@go tool cover -html=$(PROG_NAME).coverprofile ./pkg/...

lint: install-deps-dev
	@errors=$$(gofmt -l .); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi
	@errors=$$(glide novendor | xargs -n 1 golint -min_confidence=0.3); if [ "$${errors}" != "" ]; then echo "$${errors}"; exit 1; fi

vet:
	@go vet $$(glide novendor)