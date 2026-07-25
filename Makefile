.PHONY: build clean stop start restart test openspec-status openspec-validate openspec-list

build:
	go build -o todo-srv ./cmd/srv

clean:
	rm -f todo-srv

test:
	go test ./...

openspec-status:
	@openspec status || true

openspec-validate:
	@openspec validate --all

openspec-list:
	@echo "Changes:" && openspec list && echo && echo "Specs:" && openspec list --specs
