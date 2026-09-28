.PHONY: test build clean

test:
	go test ./... -coverprofile=coverage.out

build:
	go build -o telesieve ./cmd/telesieve

clean:
	rm -f telesieve coverage.out
