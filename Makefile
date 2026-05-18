.PHONY: build test test-race lint tidy vet clean release release-dry-run release-with-notes

build:
	go build -o bin/fxlint ./cmd/fxlint

test:
	go test ./...

test-race:
	go test -race ./...

vet:
	go vet ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -rf bin/ coverage.*

##@ Release

release: ## Создать новый релиз (интерактивно)
	@chmod +x dev/tools/release.sh
	@./dev/tools/release.sh

release-dry-run: ## Предпросмотр релиза без изменений
	@chmod +x dev/tools/release.sh
	@./dev/tools/release.sh --dry-run

release-with-notes: ## Создать релиз с AI-генерированными заметками (требует claude CLI)
	@chmod +x dev/tools/release.sh
	@./dev/tools/release.sh --generate-notes
