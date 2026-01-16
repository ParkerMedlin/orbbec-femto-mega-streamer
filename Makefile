.PHONY: build-windows test-windows tidy

# Windows build (default target for now)
build-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build ./...

test-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go test ./...

tidy:
	go mod tidy
