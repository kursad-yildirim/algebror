# Check if APP_NAME environment variable is set
ifndef APP_NAME
$(error APP_NAME is not set. Please set the environment variable APP_NAME)
endif

# Variables
CONTAINER_NAME := $(APP_NAME)
MAIN_GO := ./main.go
BINARY_NAME := $(APP_NAME)

# Detect the operating system
ifeq ($(OS),Windows_NT)
    # Windows
    $(error Windows is not supported. Only linux and mac!)
endif

# Clean target
clean:
	# Stop any running process with the name $(APP_NAME)
	-pkill -f $(APP_NAME)
	# Remove the binary
	rm -f $(BINARY_NAME)

# Linux build target
linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) $(MAIN_GO)

# macOS build target
macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-macos $(MAIN_GO)

# Podman build target
podman:
	podman build -t $(APP_NAME) --build-arg APP_NAME=$(APP_NAME) .

# Docker build target
docker:
	docker build -t $(APP_NAME) --build-arg APP_NAME=$(APP_NAME) .

# Clean targets for Linux, macOS, and Windows
clean-linux:
	rm -f $(BINARY_NAME)

clean-macos:
	rm -f $(BINARY_NAME)-macos


# Clean target for Podman images
clean-podman:
	# Stop and remove any running container with the app name
	-podman ps -a --filter "name=$(CONTAINER_NAME)" --format "{{.ID}}" | xargs -r podman stop
	-podman ps -a --filter "name=$(CONTAINER_NAME)" --format "{{.ID}}" | xargs -r podman rm
	# Remove dangling images
	-podman images -q --filter "dangling=true" | xargs -r podman rmi
	# Remove the image
	-podman images -q --filter "reference=$(APP_NAME)" | xargs -r podman rmi

# Clean target for Docker images
clean-docker:
	# Stop and remove any running container with the app name
	-docker ps -a --filter "name=$(CONTAINER_NAME)" --format "{{.ID}}" | xargs -r docker stop
	-docker ps -a --filter "name=$(CONTAINER_NAME)" --format "{{.ID}}" | xargs -r docker rm
	# Remove dangling images
	-docker images -q --filter "dangling=true" | xargs -r docker rmi
	# Remove the image
	-docker images -q --filter "reference=$(APP_NAME)" | xargs -r docker rmi

.PHONY: build clean linux macos windows clean-linux clean-macos clean-windows clean-podman clean-docker
