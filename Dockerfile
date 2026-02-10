FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk update && apk add make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags="-s -w" -o build/raspadmin_amd64 main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/raspadmin_amd64 /app/raspadmin

EXPOSE 8080

RUN chmod +x /app/raspadmin

CMD ["/app/raspadmin"]