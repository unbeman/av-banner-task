include ./test-app.env
export

dc-up:
	docker compose -f docker-compose.yml up

dc-up-d:
	docker compose -f docker-compose.yml up -d

dc-stop:
	docker compose -f docker-compose.yml down

dc-down:
	docker compose -f docker-compose.yml down

swagger:
	swag init -g cmd/main.go

integration-test:
	docker compose -f docker-compose-api-test.yml up -d
	docker exec -i test-banner-database psql -U postgres -d banner-keeper -a < api_test_pg_script/fill_tables.sql
	go test -tags=integration ./... -v
	docker compose -f docker-compose-api-test.yml down
