.PHONY: gqlgen up build

dc = docker-compose -f docker-compose.yaml
img_name = eu.gcr.io/vediagames/vg_api
version = latest
include .env
export $(shell sed 's/=.*//' .env)

PATH := $(PATH):$(GOPATH)/bin

gqlgen/%:
	cd $* && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate

dev:
	echo "Starting dev environment in docker"
	docker-compose up -d

down:
	echo "Shutting down the dev environment"
	docker-compose down

generate:
	echo "Regenerating code"
	go generate ./...

terraform/plan:
	cd deploy/terraform && terraform plan -out

