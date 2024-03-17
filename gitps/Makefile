.PHONY: gqlgen up build

include .env
export $(shell sed 's/=.*//' .env)

PATH := $(PATH):$(GOPATH)/bin

terraform/gcp/%:
	cd gcp && terraform $*

terraform/circleci/%:
	cd circleci && terraform $*
