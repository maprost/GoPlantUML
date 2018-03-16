FOLDERS=$(shell go list -f '{{.Dir}}' ./... | grep -v /vendor/ | grep -v /internal/testdata/)

dep:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	glide install

lint:
	gometalinter --vendor $(FOLDERS)

test:
	go test -cover $(FOLDERS)

build: dep lint test
	go build

install: dep lint test
	go install
