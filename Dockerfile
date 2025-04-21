FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download



RUN go build -o auth_service ./internal/services/auth_service/cmd

FROM alpine:latest

WORKDIR /root/
ARG PORT
ARG HOST
ENV HOST=${AUTH_SERVICE_HOST}
ENV PORT=${AUTH_SERVICE_PORT}
COPY --from=builder /app/auth_service .

COPY --from=builder /app/configs ./configs

EXPOSE ${PORT}

CMD ["./auth_service"]