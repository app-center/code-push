PROJECT := Code Push

fmt:
	@gofmt -s -w ./{cmd,daemon,gateway,pkg}/

vet:
	@go vet ./{cmd,daemon,gateway,pkg}/...

test:
	go test -count=1 -v -race ./pkg/...

benchmark:
	go test -count=1 -cpu 1 -bench . ./pkg/...