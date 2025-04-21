unit-tests:
	go test -v ./internal/... -covermode=atomic -coverprofile=coverage.out

load-test:
	k6 run hacks/k6/script.js
