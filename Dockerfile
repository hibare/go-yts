ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION:-1.24}-alpine AS base

# Build main app
FROM base AS build

WORKDIR /src/

# Add only necessary build dependencies for CGO and SQLite
RUN apk --update add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

COPY . /src/

# Enable CGO and build
RUN CGO_ENABLED=1 go build -o /bin/go_yts ./cmd/yts/main.go

# Generate final image
FROM alpine:latest

# Install CA certificates and SQLite runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    sqlite-libs

COPY --from=build /bin/go_yts /bin/go_yts

ENTRYPOINT ["/bin/go_yts"]