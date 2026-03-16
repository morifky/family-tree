.PHONY: build run stop logs clean

IMAGE_NAME ?= brayat-app
CONTAINER_NAME ?= brayat-container
PORT ?= 8080

# Build the docker image
build:
	docker build -t $(IMAGE_NAME) .

# Run the docker container in detached mode
run: stop
	docker run -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		-v brayat_data:/data \
		$(IMAGE_NAME)

# Stop and remove the container
stop:
	docker stop $(CONTAINER_NAME) 2>/dev/null || true
	docker rm $(CONTAINER_NAME) 2>/dev/null || true

# View container logs
logs:
	docker logs -f $(CONTAINER_NAME)

# Clean up container, image, and volume
clean: stop
	docker rmi $(IMAGE_NAME) 2>/dev/null || true
	docker volume rm brayat_data 2>/dev/null || true
