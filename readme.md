# GlobeChat

A place where converations happens on the world map

## Prerequists

1. Golang.
2. `make` command.
3. Postgresql
4. Docker (optional)

## How to run locally

1. Copy `.envrc.example` to `.envrc`
2. Fill `.envrc`
3. Run `make migrate_up`
4. Run `make run_backend` and `make run_frontend`

## How to build

1. Run `make`

## How to run using docker

1. Copy `.env.docker.example` to `.env.docker`
2. Fill `.env.docker`
3. Run `docker compose up`
