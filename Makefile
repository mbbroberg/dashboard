all: build test

deps:
	godep save github.com/mbbroberg/dashboard/...

build: deps
	godep go install github.com/mbbroberg/dashboard/...

test: deps
	godep go test . ./cmd/... ./triage/...
