run:
	@docker stop image-service || true && docker rm image-service || true # avoid errors in output.
	docker build -t image-service .
	docker run -d --name=image-service -p 8080:8080 image-service

start:
	docker start image-service

stop:
	docker stop image-service

logs:
	docker logs image-service

# debug:

server:
	go run cmd/server/main.go

worker:
	go run cmd/worker/main.go