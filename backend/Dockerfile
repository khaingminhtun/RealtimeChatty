# STAGE 1: Build Stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o backendapp ./cmd

# STAGE 2: The Prodcution App Image Runner
FROM scratch AS release

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/backendapp /backendapp

EXPOSE 8080

ENTRYPOINT ["/backendapp"]

# STAGE 3: The Dedicated  Migration Image
FROM golang:1.26-alpine AS migration

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENV PATH="/root/go/bin:${PATH}"

WORKDIR /migrations

COPY --from=builder /app/db/migrations .

ENTRYPOINT ["sh", "-c", "goose postgres \"$DATABASE_URL\" up -v"]