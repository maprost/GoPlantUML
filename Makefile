FULL_FOLDER_PATH=$(shell go list -f '{{.Dir}}' ./... | grep -v /vendor/ | grep -v /internal/testdata/)
FOLDER_PATH=$(shell go list ./... | grep -v /vendor/ | grep -v /internal/testdata/)

dep:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	glide install

lint:
	gometalinter --vendor $(FULL_FOLDER_PATH)

test: lint
	go test -cover $(FOLDER_PATH)

build: dep test
	go build

install: dep test
	go install
