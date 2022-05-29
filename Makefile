build: generate-mocks compile-go

test: generate-mocks run-ginkgo

ci-build: test

generate-mocks:
	mockgen -source=client/ynab/client.go -destination=client/mock_ynab/mock_client.go

run-ginkgo:
	ginkgo run -r ./..

compile-go:
	go build

release:
	env GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/autonabber
	env GOOS=windows GOARCH=amd64 go build -o dist/windows/amd64/autonabber.exe