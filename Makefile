.PHONY: gqlgen up build

dc = docker-compose -f docker-compose.yml
env_file = ./.env

include $(env_file)
export $(shell sed 's/=.*//' $(env_file))
PATH := $(PATH):$(GOPATH)/bin

gqlgen/%:
	go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate --config=$*/graphql/gqlgen.yml

dev:
	echo "Starting dev environment in docker"
	$(dc) up -d

down:
	echo "Shutting down the dev environment"
	docker-compose down

generate:
	echo "Regenerating code"
	go get github.com/99designs/gqlgen && go generate ./...

build:
	@docker build -f ./build/Dockerfile -t $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_VERSION) \
		--build-arg GITHUB_USERNAME=$(GITHUB_USERNAME) \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN) .
push:
	@docker push $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_VERSION)
