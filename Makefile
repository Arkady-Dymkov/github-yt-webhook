.PHONY: build test clean docker-build docker-run

# Build the application
build:
	go build -o github-webhook-youtrack

run:
	./github-webhook-youtrack

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f github-webhook-youtrack

# Build Docker image
docker-build:
	docker build -t github-webhook-youtrack .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e YOUTRACK_TEST_URL=${YOUTRACK_TEST_URL} \
		-e YOUTRACK_TEST_TOKEN=${YOUTRACK_TEST_TOKEN} \
		github-webhook-youtrack

# Run with docker-compose
docker-up:
	docker-compose up -d

# Stop docker-compose
docker-down:
	docker-compose down