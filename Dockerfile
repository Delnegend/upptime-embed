FROM node:lts-alpine AS node-build

WORKDIR /app

COPY ./frontend .

RUN npm install -g pnpm
RUN pnpm install --frozen-lockfile

RUN pnpm generate

FROM golang:alpine AS go-build

WORKDIR /app

COPY routes routes
COPY utils utils
COPY main.go go.mod go.sum .
COPY --from=node-build /app/.output/public ./frontend/.output/public

RUN go build -o main .

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

COPY --from=go-build /app/main .

EXPOSE 3001

CMD ["./main"]