default:
    @just --choose

build-docker:
    docker build -t Delnegend/upptime-embed .

_build_client:
    cd frontend && pnpm nuxt generate

lint:
    go fmt && cd frontend && \
        bun run oxlint --import-plugin -D correctness -D perf \
        --ignore-pattern src/dev-dist/**/*.* \
        --ignore-pattern src/utils/artefact-wasm/**/*.* && \
        bun run prettier -l -w "**/*.{js,ts,vue,json,css}"

build-static: _build_client
    go build -o main .

dev-server:
    go run .

dev-client:
    cd frontend && bun run nuxt dev