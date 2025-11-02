# upptime-api

Custom upptime frontend page and APIs without rate limits.

> Regex ftw.

## Usage
```
https://upptime.delnegend.com/?user=upptime&repo=upptime
```

Replace `upptime` with your GitHub username and the upptime repo name.

## Self-hosting
1. Clone the repo
    ```bash
    git clone https://github.com/Delnegend/upptime-embed.git
    cd upptime-embed
    ```

2. Copy `docker-compose.example.yml` to `docker-compose.yml` and edit the ports if needed.
    ```bash
    cp docker-compose.example.yml docker-compose.yml
    ```

3. Start the services
    ```bash
    docker compose up -d
    ```

## Development
- Client:
    ```bash
    cd frontend
    pnpm i && pnpm dev
    ```

- Server
    ```bash
    go run main.go
    ```

- Caddy
    ```bash
    caddy run --config Caddyfile
    ```

## Exposed endpoints

### `API` `GET /api/alive`

Returns 200 OK

### `API` `GET /api/{username}/{reponame}`

Returns JSON
```json
{
    "Overall": "all_good" | "degraded" | "down" | "partial" | "unknown",
    "Details": [
        {
            "iconUrl": string,
            "title": string,
            "slug": string,
            "url": string | null,
            "status": "up" | "down" | "degrade" | "unknown",
            "responseOverall": string,
            "response24h": string,
            "response7d": string,
            "response30d": string,
            "response1y": string,
            "uptimeOverall": string,
            "uptime24h": string,
            "uptime7d": string,
            "uptime30d": string,
            "uptime1y": string,
        }
    ]
}
```

### `API` `GET /api/graph/{username}/{reponame}/{slug}/{duration}`

Valid durations: `day`, `week`, `month`, `year`, `all`

Returns PNG

### `FRONTEND` `GET /?user={username}&repo={reponame}`

## License

MIT