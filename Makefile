run:
	go run ./cmd/api

db-up:
	podman-compose up -d

db-start:
	podman-compose start
