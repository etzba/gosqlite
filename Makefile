NAME ?= gosqlite
TAG ?= latest
REPO ?= etzba/${NAME}

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

lint:
	golangci-lint run ./...

# test with etzba 
exec:
	etz api --exec=scripts/exec.yaml -d=3s -w=2 -r=10 --output=results/$$(date +%Y%m%d_%H%M%S)_result.json

# build \ push to dockerhub
.PHONY: docker-build
docker-build:
	docker build -t ${REPO}:${TAG} .

.PHONY: docker-push
docker-push:
	docker push ${REPO}:${TAG}

.PHONY: cleanup
cleanup:
	docker rm $$(docker stop $$(docker ps -a -q --filter ancestor=${REPO}:${TAG} --format="{{.ID}}"))

.PHONY: run
run:
	docker run -e RELEASE_STAGE=testing -p 127.0.0.1:8080:8080/tcp -t --rm ${REPO}:${TAG}
