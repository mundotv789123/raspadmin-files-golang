FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
RUN go build -ldflags="-s -w" -o build/raspadmin main.go

FROM alpine:3.23

RUN apk update && apk add ffmpeg ffmpegthumbnailer

ENV GIN_MODE=release

ENV DB_FILE=data/database.db
RUN mkdir -p /app/data

WORKDIR /app

COPY --from=builder /app/build/raspadmin /usr/local/bin/raspadmin

EXPOSE 8080

RUN chmod +x /usr/local/bin/raspadmin

CMD ["raspadmin"]