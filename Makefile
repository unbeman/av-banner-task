DB_URL=postgresql://postgres:password@localhost:6500/banner-keeper?sslmode=disable


create-db:
	docker exec -it banner-database createdb --username=postgres --owner=postgres banner-keeper

drop-db:
	docker exec -it banner-database dropdb banner-keeper

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down