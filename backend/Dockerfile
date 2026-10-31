# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tini \
    && addgroup -S app \
    && adduser -S app -G app
WORKDIR /app
COPY --from=builder /server /app/server
COPY migrations /migrations
ENV MIGRATIONS_DIR=/migrations \
    DB_PATH=/data/taxilux.db \
    PORT=8080 \
    LOG_LEVEL=info
RUN mkdir -p /data && chown app:app /data
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost:8080/api/health || exit 1
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/app/server"]
