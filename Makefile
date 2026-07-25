.PHONY: build clean stop start restart test test-js openspec-status openspec-validate openspec-list

build:
	go build -o todo-srv ./cmd/srv

clean:
	rm -f todo-srv

test:
	go test ./...

test-js:
	node srv/static/script.test.js

openspec-status:
	@openspec status || true

openspec-validate:
	@openspec validate --all

openspec-list:
	@echo "Changes:" && openspec list && echo && echo "Specs:" && openspec list --specs
