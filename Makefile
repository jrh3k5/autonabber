build: generate-mocks compile-go

test: generate-mocks run-ginkgo

ci-build: test

generate-mocks:
	mockgen -source=client/ynab/client.go -destination=client/mock_ynab/mock_client.go

run-ginkgo:
	ginkgo run -r ./..

compile-go:
	go build