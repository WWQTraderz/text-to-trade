.PHONY: build-images
build-images:
	docker buildx build -f Dockerfile.migrations -t tylerschade268/ttt-migrations:latest .
	docker buildx build -f Dockerfile.app -t tylerschade268/ttt-api:latest .

push-images:
	docker push tylerschade268/ttt-migrations:latest
	docker push tylerschade268/ttt-api:latest