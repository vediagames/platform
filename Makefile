.PHONY: gqlgen up build

dc = docker-compose -f docker-compose.yaml
img_name = eu.gcr.io/vediagames/vg_api
version = latest
#env_file = ./.env

#include $(env_file)
export $(shell sed 's/=.*//' $(env_file))
PATH := $(PATH):$(GOPATH)/bin

gqlgen/%:
	cd $* && go get github.com/99designs/gqlgen && go run github.com/99designs/gqlgen generate

postgres:
	docker run -e "POSTGRES_USER=vedia" -e "POSTGRES_PASSWORD=123" -e "POSTGRES_DB=vediagames" -p 5432:5432 -d postgres:15.0-alpine

redis:
	docker run -p 6379:6379 -d redis