IMAGE_NAME ?= "osodevops/beats-aws-ecs"
GO_VERSION ?= "1.12.2"
GO_PLATFORM ?= "linux-amd64"
BEATS_VERSION ?= "7.5.0"
BEAT_NAME ?= "filebeat"
BEAT_DOCKER_IMAGE = "docker.elastic.co/beats/${BEAT_NAME}:${BEATS_VERSION}"
BEATS_AWS_ECS_VERSION ?= "0.1.0"

build:
	docker build \
		-t $(IMAGE_NAME) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg GO_PLATFORM=$(GO_PLATFORM) \
		--build-arg BEATS_VERSION=$(BEATS_VERSION) \
		--build-arg BEAT_NAME=$(BEAT_NAME) \
		--build-arg BEATS_AWS_ECS_VERSION=$(BEATS_AWS_ECS_VERSION) \
		--build-arg BEAT_DOCKER_IMAGE=$(BEAT_DOCKER_IMAGE) \
		.
