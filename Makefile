export DOCKER_BUILDKIT = 1

.PHONY: build-images
build-images:
	docker buildx build --push --platform linux/amd64,linux/arm64 -f Dockerfile.migrations -t tylerschade268/ttt-migrations:v1 .
	docker buildx build --push --platform linux/amd64,linux/arm64 -f Dockerfile.app -t tylerschade268/ttt-api:v1 .

push-images:
	docker push tylerschade268/ttt-migrations:v1
	docker push tylerschade268/ttt-api:v1

build-migrations:
	docker buildx build --push --platform linux/amd64 -f Dockerfile.migrations -t tylerschade268/ttt-migrations:v1.2 .

push-migrations:
	docker push tylerschade268/ttt-migrations:v1.2

new-migrations: build-migrations push-migrations