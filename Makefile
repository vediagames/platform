.PHONY: gqlgen up build

dc = docker-compose -f docker-compose.yaml
img_name = eu.gcr.io/vediagames/vg_api
version = latest

PATH := $(PATH):$(GOPATH)/bin

gqlgen/%:
	cd $* && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate

dev:
	echo "Starting database environment in docker"
	docker-compose up -d

down:
	echo "Shutting the environment down"
	docker-compose down

generate:
	echo "Regenerating graphql schemas"
	go generate ./...