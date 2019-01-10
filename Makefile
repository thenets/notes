IMAGE_TAG=thenets/notes

# Build Docker image
build:
	docker build -t $(IMAGE_TAG) .

# Run server in the development mode
run: 
	docker run -it --rm \
		-v $(PWD)/src:/app \
		-p 5000:5000 \
		$(IMAGE_TAG)

# Run shell inside the container
shell:
	docker run -it --rm \
		-v $(PWD)/src:/app \
		$(IMAGE_TAG)