# This Makefile was done using 'domake'
# Generated at 26/08/2025


# =================================================================================== #
# HELPERS
# =================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# =================================================================================== #
# COMMANDS
# =================================================================================== #

## build: Builds the app
.PHONY: build
build: 
	@echo ":: Beginning to build the app..."
	@go build -o mangathorg -ldflags="-w -s" -trimpath ./cmd/
	

## dev: Runs the app in the development environment
.PHONY: dev
dev: 
	@echo ":: Running the app in dev environment..."
	@go run ./cmd/
	

## docker/build: Builds and runs the app with docker in interactive mode
.PHONY: docker/build
docker/build: 
	@echo ":: Beginning to build the docker image..."
	@docker compose --profile dev up --build
	

## docker/deploy: Deploys the last version of the docker container to the self-hosted docker repository
.PHONY: docker/deploy
docker/deploy:  confirm
	@echo ":: Building and pushing the docker image..."
	@docker compose build mangathorg-dev && \
	docker tag mangathorg-mangathorg-dev docker.adebarbarin.com/mangathorg && \
	docker push docker.adebarbarin.com/mangathorg && \
	echo ":: Job finished successfully!" || (echo ":: Job failed." && exit 1)

