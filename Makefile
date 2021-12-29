cat := $(if $(filter $(OS),Windows_NT),type,cat)
VERSION := $(shell $(cat) VERSION)

build:
	CGO_ENABLED=0 GOARCH=amd64 go build -o swarm-cronjob -a -installsuffix cgo --ldflags="-s -w -X 'main.version=${VERSION}'" ./cmd/main.go

build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o swarm-cronjob -a -installsuffix cgo --ldflags="-s -w -X 'main.version=${VERSION}'" ./cmd/main.go

docker-build: build-linux
	docker build -t devops/swarm-cronjob-event:${VERSION} -f ./Dockerfile.event .

publish:
	docker tag devops/swarm-cronjob-event:${VERSION} 10.8.32.17:8082/devops/swarm-cronjob-event:${VERSION}
	docker push 10.8.32.17:8082/devops/swarm-cronjob-event:${VERSION}

run:
	./swarm-cronjob