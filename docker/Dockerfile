FROM golang:latest

WORKDIR /app

COPY ./ /app

RUN /bin/sh /app/docker/tuner.sh

CMD ["/app/bin/app"]