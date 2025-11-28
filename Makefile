build: generate-mocks compile-go

test: generate-mocks run-ginkgo

ci-build: test build

generate-mocks:
	mockgen -source=client/ynab/client.go -destination=client/mock_ynab/mock_client.go

run-ginkgo:
	go test ./...

compile-go:
	go build

release-clean:
	rm -rf dist

release-build-mac-x64:
	env GOOS=darwin GOARCH=amd64 go build -o dist/darwin/amd64/autonabber
	tar -C dist/darwin/amd64/ -czvf dist/darwin/amd64/osx-x64.tar.gz autonabber

release-build-mac-arm64:
	env GOOS=darwin GOARCH=arm64 go build -o dist/darwin/arm64/autonabber
	tar -C dist/darwin/arm64/ -czvf dist/darwin/arm64/osx-arm64.tar.gz autonabber

release-build-win-x64:
	env GOOS=windows GOARCH=amd64 go build -o dist/windows/amd64/autonabber.exe
	(cd dist/windows/amd64 && zip -r - autonabber.exe) > dist/windows/amd64/win-x64.zip

release-build: release-build-mac-x64 release-build-mac-arm64 release-build-win-x64

release: release-clean release-build