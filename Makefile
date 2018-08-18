GO_PKGS = $(shell go list ./... | grep -v /vendor/)

vet:
	go vet $(GO_PKGS)

test:
	go test $(GO_PKGS)

bench:
	go test -run=XXX -bench .
