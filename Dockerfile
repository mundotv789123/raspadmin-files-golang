FROM alpine:3.23

RUN apk update && apk add ffmpeg ffmpegthumbnailer

ENV GIN_MODE=release

ENV DB_FILE=data/database.db
RUN mkdir -p /app/data

WORKDIR /app

COPY ./build/raspadmin /usr/local/bin/raspadmin

EXPOSE 8080

RUN chmod +x /usr/local/bin/raspadmin

CMD ["raspadmin"]