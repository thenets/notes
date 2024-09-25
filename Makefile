IMAGE_TAG=thenets/notes
CONTAINER_RUNTIME=podman

container-build:
	${CONTAINER_RUNTIME} build -t $(IMAGE_TAG) .

container-run:
	${CONTAINER_RUNTIME} run -it --rm \
		-v $(PWD)/src:/app \
		-p 8080:8080 \
		$(IMAGE_TAG)

container-shell:
	${CONTAINER_RUNTIME} run -it --rm \
		--name thenets-notes \
		-v $(PWD)/src:/app \
		$(IMAGE_TAG) /bin/sh
