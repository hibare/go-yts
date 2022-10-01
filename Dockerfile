FROM golang:1.19.1-alpine AS base

# Build main app
FROM base AS build

WORKDIR /src/

RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
    update-ca-certificates

COPY . /src/

RUN CGO_ENABLED=0 go build -o /bin/go_yts

# Generate final image
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/go_yts /bin/go_yts

ENTRYPOINT ["/bin/go_yts"]