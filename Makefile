unit-tests:
	go test -v ./internal/... -covermode=atomic -coverprofile=coverage.out
