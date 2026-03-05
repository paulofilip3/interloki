.PHONY: build build-frontend build-backend dev clean lint

build: build-frontend build-backend

build-frontend:
	cd web && npm run build

build-backend:
	go build -o interloki ./cmd/interloki

dev:
	go run ./cmd/interloki

clean:
	rm -f interloki
	rm -rf web/dist

lint:
	go vet ./...
	golangci-lint run ./...
