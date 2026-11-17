run:
	go run ./cmd/api

db-up:
	podman-compose up -d
