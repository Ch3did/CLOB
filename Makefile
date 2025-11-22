APP_NAME=clob-api
DOCKER_COMPOSE=docker-compose

.PHONY: help up down restart logs build api test clean

help:
	@echo ""
	@echo "Comandos disponíveis:"
	@echo "  make up         - sobe os containers (API + Postgres)"
	@echo "  make down       - derruba os containers"
	@echo "  make restart    - reinicia tudo"
	@echo "  make logs       - mostra os logs da API"
	@echo "  make build      - build da aplicação"
	@echo "  make clean      - remove containers/volumes"

# ----------------------------
# DOCKER
# ----------------------------

up:
	$(DOCKER_COMPOSE) up --build -d

down:
	$(DOCKER_COMPOSE) down

restart: down up

logs:
	$(DOCKER_COMPOSE) logs -f api

clean:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

# ----------------------------
# GO LOCAL
# ----------------------------

build:
	go build -o $(APP_NAME) ./api/main.go
