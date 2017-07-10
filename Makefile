all: bin

TAG = 0.4.1b
NAME = journald2graylog

PREFIX = hub.docker.com
ifneq ($(strip $(DOCKER_REGISTRY)),)
	PREFIX = $(DOCKER_REGISTRY)
endif
REGISTRY := $(PREFIX)/$(NAME)

GORUN = go run
GOBUILD = go build -v -a
GOBUILD_RELEASE = go build -a -ldflags '-s -w'

run:
	cat fixtures/testset-4l.jsonl | J2G_HOSTNAME=localhost $(GORUN) journald2graylog.go -v

bin_linux: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(NAME)

bin_mac: clean
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(NAME)

bin: clean
	 $(GOBUILD) -o $(NAME)

docker_image: bin_linux
	@docker build \
		--build-arg VCS_REF=`git rev-parse --short HEAD` \
		--build-arg BUILD_DATE=`date +%Y%m%d` \
		--build-arg BUILD_HOST=`hostname` \
		-t $(REGISTRY):$(TAG) .

clean:
	rm -rf $(NAME) dist/

docker_push: docker_image
	docker push $(REGISTRY):$(TAG)

dist: clean
	mkdir -p dist/darwin
	GOOS=darwin GOARCH=amd64 $(GOBUILD_RELEASE) -o dist/darwin/$(NAME)
	cd dist/darwin && tar -czvf $(NAME)-v$(TAG).tgz $(NAME)
	rm -rf dist/darwin/$(NAME)

	mkdir -p dist/linux
	GOOS=linux GOARCH=amd64 $(GOBUILD_RELEASE) -o dist/linux/$(NAME)
	cd dist/linux && tar -czvf $(NAME)-v$(TAG).tgz $(NAME)
	rm -rf dist/linux/$(NAME)

	help:
		$(info Usage:)
		$(info use: "make" to build the binary for the current platform)
		$(info use: "make docker_image" to build a Linux binary and it's docker image)
		$(info use: "make docker_push" to push the docker image to hub.docker.com)
