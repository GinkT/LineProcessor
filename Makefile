lint:
	go vet ./...

tests:
	go test ./...

run:
	docker-compose up

stop:
	docker-compose stop
