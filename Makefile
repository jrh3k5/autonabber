build: generate-mocks compile-go

test: generate-mocks run-ginkgo

ci-build: test build

generate-mocks:
	mockgen -source=client/ynab/client.go -destination=client/mock_ynab/mock_client.go

run-ginkgo:
	ginkgo run -r ./..

compile-go:
	go build

release-clean:
	rm -rf dist

release-build:
	env GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/autonabber
	tar -C dist/darwin/amd64/ -czvf dist/darwin/amd64/osx-x64.tar.gz autonabber
	env GOOS=windows GOARCH=amd64 go build -o dist/windows/amd64/autonabber.exe
	(cd dist/windows/amd64 && zip -r - autonabber.exe) > dist/windows/amd64/win-x64.zip

release: test release-clean release-build