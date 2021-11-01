generate:
	go generate ./...

docker-compose-up:
	docker-compose up --detach

docker-compose-down:
	docker-compose down