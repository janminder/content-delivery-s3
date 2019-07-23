GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
GO_DEV_ENV := CGO_ENABLED=0 GOOS=darwin GOARCH=amd64

DOCKER_BUILD=$(shell pwd)/bin
DOCKER_CMD=$(DOCKER_BUILD)/cd-s3

$(DOCKER_CMD): clean
	mkdir -p $(DOCKER_BUILD)
	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) ./cmd/...

run: clean
	$(GO_DEV_ENV) go build -v -o $(DOCKER_CMD) ./cmd/...
	exec $(DOCKER_CMD)

clean:
	rm -rf $(DOCKER_BUILD)

get:
	dep ensure

push: clean
	cf push
