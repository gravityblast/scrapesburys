.PHONY: test

GO=GO15VENDOREXPERIMENT=1 go

test:
	$(GO) test -v
