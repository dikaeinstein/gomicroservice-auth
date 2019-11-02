FROM datadog/dogstatsd:6.15.0-rc.8

RUN mkdir /service
WORKDIR /service

RUN apk update && apk add ca-certificates \
    && apk add supervisor \
    && mkdir -p /var/log/supervisor

COPY ./auth /service/auth
COPY ./supervisor.conf /service/supervisor.conf

EXPOSE 8081

CMD ["supervisord", "-c", "/service/supervisor.conf"]
