PROJECT := Code Push

fmt:
	@gofmt -s -w ./{cmd,daemon,gateway,pkg}/

vet:
	@go vet ./{cmd,daemon,gateway,pkg}/...

test:
	go test -count=1 ./pkg/...

test-version-compat-tree:
	go test ./daemon/code-push/usecase/version_compat_tree/

benchmark:
	go test -count=1 -cpu 1 -bench . ./pkg/...

generate:
	@go generate ./...

@phony: fmt vet test benchmark generate