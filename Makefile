all: bin

TAG = 0.3.0b
NAME = journald2graylog

PREFIX = hub.docker.com
ifneq ($(strip $(DOCKER_REGISTRY)),)
	PREFIX = $(DOCKER_REGISTRY)
endif
REGISTRY := $(PREFIX)/$(NAME)

GOBUILD = go build -v -a
GOBUILD_RELEASE = go build -a -ldflags '-s -w'

bin_linux: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(NAME)

bin_mac: clean
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(NAME)

bin: clean
	 $(GOBUILD) -o $(NAME)

docker: bin_linux
	docker build --pull -t $(REGISTRY):$(TAG) .

clean:
	rm -rf $(NAME) dist/

push: container
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
