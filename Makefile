all: push

# 0.0 shouldn't clobber any release builds, current "latest" is 0.9
TAG = 0.1.3

PREFIX = hub.docker.com
ifneq ($(strip $(DOCKER_REGISTRY)),)
	PREFIX = $(DOCKER_REGISTRY)
endif
REGISTRY := $(PREFIX)/journald2graylog

controller: clean
	GOOS=linux go build -a -ldflags '-w' -o journald2graylog

container: controller
	docker build --pull -t $(PREFIX):$(TAG) .

clean:
	rm -f journald2graylog

push: container
	#gcloud docker push $(PREFIX):$(TAG)
