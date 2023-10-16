
run:
	go run cmd/app/main.go
generate-mocks:
	mockery --all --keeptree
test-verbose:
	go test -v --cover ./...
test:
	go test -v -covermode=count -coverprofile=coverage.out $(shell go list ./... | egrep -v '/mocks|/constant|/entity|/model') -json > report.json
coverage:
	go tool cover -func=coverage.out
coverage-html:
	go tool cover -html=coverage.out