default: lint test coverage clean

clean: # Remove temporary files
	rm -f ./coverage.out

coverage: # Generate coverage badge
	.build/coverage.sh

codegen:
	protoc \
		--go_out=. \
		--go-grpc_out=. \
		--proto_path=services/protocol \
		port.proto

lint: # Analyse code cleanliness
	docker run \
		--rm \
		--volume $(PWD):/src \
		--workdir /src \
		golangci/golangci-lint \
		golangci-lint run -v

test: # Run tests
	docker-compose \
		--file .build/docker-compose.yml \
		up \
		--build \
		--force-recreate \
		--abort-on-container-exit
