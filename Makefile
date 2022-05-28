build:
	go build

test:
	ginkgo run -r ./..

ci-build:
	test