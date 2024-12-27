FROM golang:1.23-alpine AS builder

RUN apk update && apk add
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o app ./cmd/order/main.go

FROM alpine:3.20

WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8000
CMD ["./app"]