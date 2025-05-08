up:
	docker-compose -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down -v

migrate:
	docker run --rm --network sql_default -v ./flyway:/flyway/project flyway/flyway \
	-workingDirectory="project" \
	-configFiles=flyway.toml \
	migrate

baseline:
	docker run --rm --network sql_default -v ./flyway:/flyway/project flyway/flyway \
	-workingDirectory="project" \
	-configFiles=flyway.toml \
	baseline

repair:
	docker run --rm --network sql_default -v ./flyway:/flyway/project flyway/flyway \
	-workingDirectory="project" \
	-configFiles=flyway.toml \
	repair

clean:
	docker run --rm --network sql_default -v ./flyway:/flyway/project flyway/flyway \
	-workingDirectory="project" \
	-configFiles=flyway.toml \
	clean