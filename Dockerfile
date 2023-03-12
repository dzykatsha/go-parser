FROM golang:latest

RUN apt-get update -qq
RUN apt-get install -y -qq antiword

