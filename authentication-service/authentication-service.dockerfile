# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY /bin/authApp /app

CMD [ "/app/authApp" ]