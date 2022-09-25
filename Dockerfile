FROM golang:1.19.0-alpine AS base

# Build main app
FROM base AS build

WORKDIR /src/

COPY . /src/

RUN CGO_ENABLED=0 go build -o /bin/go_yts

# Generate final image
FROM scratch

COPY --from=build /bin/go_yts /bin/go_yts

ENTRYPOINT ["/bin/go_yts"]