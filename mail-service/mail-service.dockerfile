# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /bin/mailerApp /app

CMD [ "/app/mailerApp" ]