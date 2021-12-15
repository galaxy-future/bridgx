vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

check: vet

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

# Quick start
# Pull images from dockerhub and run
docker-run-linux:
	sh ./run-for-linux.sh

docker-run-mac:
	sh ./run-for-mac.sh

docker-container-stop:
	docker ps -aq | xargs docker stop
	docker ps -aq | xargs docker rm

docker-image-rm:
	docker image prune --force --all

# Immersive experience
# Compile and run by docker-compose
docker-compose-start:
	docker-compose up -d

docker-compose-stop:
	docker-compose down

docker-compose-build:
	docker-compose build

docker-tag:
	docker tag bridgx_api:huawei galaxyfuture/bridgx-api:huawei
	docker tag bridgx_scheduler:huawei galaxyfuture/bridgx-scheduler:huawei

docker-push-hub:
	docker push galaxyfuture/bridgx-api:huawei
	docker push galaxyfuture/bridgx-scheduler:huawei

docker-hub-all: docker-compose-build docker-tag docker-push-hub