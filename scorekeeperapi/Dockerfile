FROM golang:1.14-buster
LABEL MAINTAINER="tank@hey.com"

RUN apt-get update; apt-get upgrade; apt-get install -y sqlite3

RUN go get -u github.com/gorilla/mux; \
    go get database/sql; \ 
    go get github.com/mattn/go-sqlite3; \
    go get github.com/stretchr/testify;

WORKDIR /code

COPY code/* /code/
RUN go build -o scorekeeperapi main.go

RUN echo 'alias ll="ls -al"' >> ~/.bashrc;

ENTRYPOINT [ "/code/scorekeeperapi" ]