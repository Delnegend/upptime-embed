default:
    @just --choose

build-docker:
    docker build -t Delnegend/upptime-embed .

_build_client:
    cd frontend && pnpm nuxt generate

lint:
    cd frontend && pnpm eslint --cache --fix .

build-static: _build_client
    go build -o main .

extract-docker-image:
    docker save Delnegend/upptime-embed | gzip > docker-image.tar.gz