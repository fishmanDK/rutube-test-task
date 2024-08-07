FROM golang:1.21.1-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

RUN apk update && apk add postgresql-client

COPY wait-for-postgres.sh /usr/local/src/wait-for-postgres.sh

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app cmd/rutubeTestTask/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /app

COPY .env .env
COPY config/config.yml config/config.yml
COPY config/config-db.yml config/config-db.yml
RUN ls -la /app
WORKDIR /

CMD ["/app"]