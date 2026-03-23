FROM golang:1.25-alpine AS base

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/server"]

FROM base AS builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/server ./cmd/server

FROM alpine:3.21 AS prod

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -H -u 10001 appuser \
	&& mkdir -p /app/configs /app/logs

COPY --from=builder /out/server /app/server
COPY configs/config.prod.yaml /app/configs/config.prod.yaml

ENV APP_CONFIG=/app/configs/config.prod.yaml

USER appuser

EXPOSE 8080

CMD ["/app/server"]
