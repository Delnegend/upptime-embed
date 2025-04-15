default:
    @just --choose

build-docker:
    docker build -t Delnegend/upptime-embed .

_build_client:
    cd frontend && pnpm nuxt generate

lint:
    cd frontend && pnpm eslint --cache --fix .

build: _build_client
    go build -o main .