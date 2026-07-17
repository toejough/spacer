.PHONY: build clean stop start restart test

build:
	go build -o todo-srv ./cmd/srv

clean:
	rm -f todo-srv

test:
	go test ./...
