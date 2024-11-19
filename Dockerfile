FROM golang:1.22.2-alpine3.19 AS builder
# Instalar git
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o app-crud-user ./cmd/main.go

FROM alpine as final
LABEL service="img-crud-user"
ARG version
ENV VERSION=1.0
ENV MYSQL_HOST=198.251.66.112
ENV MYSQL_PORT=3306
ENV MYSQL_USER=usr_tmp
ENV MYSQL_PASSWORD=123
ENV MYSQL_NAME=COLONIAL
ENV MYSQL_TIMEOUT=5s
ENV MYSQL_TIMEOUT_QUERY=5s
ENV MYSQL_MAX_CONNECTIONS=50
ENV MYSQL_MAX_IDLE_CONNECTIONS=5
ENV MYSQL_MAX_CONNECTION_LIFE_TIME=3m

ENV PORT=5002
ENV SHUTDOWN_TIMEOUT=5s

WORKDIR /app

COPY --from=builder /app/app-crud-res ./app-crud-user
EXPOSE 5002

CMD ["/app/app-crud-user"]