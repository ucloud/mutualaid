FROM golang:1.17.11-alpine3.16

RUN mkdir -p /app /data/conf
COPY ./bin/mutualaid /app
COPY ./bin/config.yaml /data/conf

WORKDIR /app

EXPOSE 28000
EXPOSE 29000

CMD ["./mutualaid", "-conf", "/data/conf", "-log", "/var/log/mutualaid/api.log"]
