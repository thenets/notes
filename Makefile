IMAGE_TAG=quay.io/thenets/notes
CONTAINER_RUNTIME=podman

# Go
.PHONY: go-update
go-update:
	go mod tidy

# Container
.PHONY: container-build
container-build: go-update
	${CONTAINER_RUNTIME} build -t $(IMAGE_TAG) .

.PHONY: container-run
container-run:
	${CONTAINER_RUNTIME} run -it --rm \
		-p 8080:8080 \
		$(IMAGE_TAG)

.PHONY: container-push
container-push:
	${CONTAINER_RUNTIME} push $(IMAGE_TAG)

.PHONY: container-shell
container-shell:
	${CONTAINER_RUNTIME} run -it --rm \
		--user root \
		--name thenets-notes \
		--workdir /app \
		-v $(PWD):/app:Z \
		docker.io/golang:latest /bin/sh

.PHONY: release
release: container-build container-push

.PHONY: clean
clean:
	rm -f notes
	rm -rf vendor/
