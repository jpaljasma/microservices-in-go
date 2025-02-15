# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /bin/loggerApp /app

CMD [ "/app/loggerApp" ]