# Go-YTS

Golang program to watch for `Popular Downloads` on YTS and send notifications.

Script uses web scrapping methodology to check for popular downloads.

## Deployment / Execution

There are two ways to run this program.

1. Run binary directly on host.
2. Run in Docker

### Run binary directly on host

For each release, binaries are published on Github release page using goreleaser.

- Download platform specific binary from Github release page.
- Build binary on the host

To build binary on the host, clone this repo and execute following command in the root of project directory.

```shell
go build -o /bin/go_yts ./cmd/yts/main.go
```

Rename file `app.env.example` to `app.env` and populate all environment variables required for the program to execute. Alternatively, you can export environment variables.

### Run in Docker

go-cyts is packaged as docker container. Docker image is available on [Docker Hub](https://hub.docker.com/r/hibare/go-yts).

Use following docker-compose.yml definition to run it in Docker.

```shell
version: "3.7"
services:
  go-yts:
    image: hibare/go-yts
    container_name: go-yts
    hostname: go-yts
    restart: always
    environment:
      - SCHEDULE=0 */4 * * *
      - DATA_DIR=/data
      - HISTORY_FILE=history.json
    volumes:
        - go-yts:/data
volumes:
  go-yts:
```

## Environment Variables

| Variable                 | Description                                         | Default Value  |
| ------------------------ | --------------------------------------------------- | -------------- |
| SCHEDULE                 | Internal cron schedule. Uses standard cron notation | 0 \/4 \* \* \* |
| DATA_DIR                 | Directory to store history file                     | /data          |
| HISTORY_FILE             | History filename                                    | history.json   |
| HTTP_REQUEST_TIMEOUT     | Request timeout value for scrapper                  | 60 Seconds     |
| NOTIFIER_DISCORD_WEBHOOK | Discord notification webhook                        | -              |
| NOTIFIER_DISCORD_ENABLED | Discord notification status                         | false          |
| LOG_LEVEL                | Log Level (INFO, ERROR, WARN, DEBUG)                | INFO           |
| LOG_MODE                 | Log mode (PRETTY, JSON)                             | PRETTY         |

## Notifications

Currently, only Discord is supported as notification destinations.
