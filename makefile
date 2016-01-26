.PHONY: test vet

GO=GO15VENDOREXPERIMENT=1 go

default: test

vet:
	$(GO) vet

test: vet
	$(GO) test -v -race
