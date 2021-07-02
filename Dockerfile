# == Create the binary. ======================================================
FROM golang:alpine AS build-env

ENV GO111MODULE=on

WORKDIR /app
ADD . /app
RUN cd /app && go mod download && go build -o goapp .
# ============================================================================



# == Run the binary in a small image. ========================================
FROM alpine

# Add ca-certificates.
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build-env /app/goapp /app

ENTRYPOINT ./goapp
# ============================================================================
