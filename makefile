.PHONY: test vet lint

GO=GO15VENDOREXPERIMENT=1 go

default: test

lint:
	ls | grep -v vendor | grep ".go" | xargs golint

vet:
	$(GO) vet

test: vet lint
	$(GO) test -v -race
