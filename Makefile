format:
	#go get golang.org/x/tools/cmd/goimports
	find . -name '*.go' | grep -Ev 'vendor|thrift_gen' | xargs goimports -w

build:
	sh ./scripts/build_api.sh && sh ./scripts/build_scheduler.sh

run:
	sh ./output/run_api.sh

clean:
	rm -rf output

server: format clean build run

docker-build-scheduler:
	docker build -t 172.16.16.172:12380/bridgx/bridgx-scheduler:v0.2 -f ./SCHEDULER.Dockerfile ./

docker-build-api:
	docker build -t 172.16.16.172:12380/bridgx/bridgx-api:v0.2 -f ./API.Dockerfile ./

docker-push-scheduler:
	docker push 172.16.16.172:12380/bridgx/bridgx-scheduler:v0.2

docker-push-api:
	docker push 172.16.16.172:12380/bridgx/bridgx-api:v0.2

docker-all: clean docker-build-scheduler docker-build-api docker-push-scheduler docker-push-api

# quick start
# pull images from dockerhub and run
docker-run-linux:
	sh ./run-for-linux.sh

docker-run-mac:
	sh ./run-for-mac.sh

docker-container-stop:
	docker ps -aq | xargs docker stop
	docker ps -aq | xargs docker rm

docker-image-rm:
	docker image prune --force --all