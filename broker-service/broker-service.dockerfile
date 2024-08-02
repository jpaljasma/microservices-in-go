# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /bin/brokerApp /app

CMD [ "/app/brokerApp" ]