run:
	sudo docker run -p 4222:4222 -p 8222:8222 --rm nats-streaming
	go run ./publisher/publisher.go
	go run ./cmd/main.go