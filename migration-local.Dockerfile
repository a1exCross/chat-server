FROM alpine:3.13

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

COPY internal/migrations/*.sql migrations/
COPY migration-local.sh .
COPY local.env .

RUN chmod +x migration-local.sh

ENTRYPOINT ["bash", "migration-local.sh"]