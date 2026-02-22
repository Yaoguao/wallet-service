FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk --no-cache add bash make gcc git musl-dev

# install dependency
COPY go.mod go.sum ./
RUN go mod download

# copy common code
COPY . .

RUN go build -o /app/bin/ ./cmd/wallet-service

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# BUILD WALLET-SERVICE ------------------------------------------------------------------------------------------------

FROM alpine AS wallet-service

COPY --from=builder /app/bin/wallet-service /wallet-service
COPY --from=builder /app/config /config

CMD ["/wallet-service"]

# BUILD MIGRATOR ------------------------------------------------------------------------------------------

FROM alpine AS migrator

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations

COPY /scripts/migrate.sh migrate.sh

RUN chmod +x /app/migrate.sh

ENTRYPOINT ["/app/migrate.sh"]

#ENTRYPOINT ["/bin/sh", "-c", "migrate -path=/app/migrations -database=$ up"]