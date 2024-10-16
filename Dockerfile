FROM golang:1.23-alpine AS builder

RUN apk update && apk add
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o health-service .

FROM alpine:3.20

WORKDIR /root/
COPY --from=builder /app/health-service .
EXPOSE 8000
CMD ["./health-service"]