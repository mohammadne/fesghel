codecov-test:
	go test -v ./internal/... -covermode=atomic -coverprofile=coverage.out
