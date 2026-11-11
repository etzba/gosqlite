NAME ?= gosqlite
TAG ?= latest
REPO ?= ghcr.io/etzba/${NAME}

all: test up exec down

# docker
up:
	docker-compose down
	docker-compose up -d gosqlite

down:
	docker-compose down 

cleanup:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=etzba/gosqlite:latest --format="{{.ID}}"))

# tests, lint and run
test:
	go test -v ./...

# TODO: set test from golang client
test-e2e:
	echo e2e

run:
	go run main.go

lint:
	golangci-lint run ./...

# test with etzba 
exec:
	etz api --auth=etzba/secret.yaml --exec=etzba/executions.yaml -d=3s -w=6 -r=24 --output=etzba/results/$$(date +%Y%m%d_%H%M%S)_result.json

# build \ push to dockerhub
.PHONY: docker-build
docker-build:
	docker build -t ${REPO}:${TAG} .

.PHONY: docker-push
docker-push:
	docker push ${REPO}:${TAG}

