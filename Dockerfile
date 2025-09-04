FROM ubuntu:22.04 AS dev

RUN apt update && apt install -y wget tar git delve vim bash
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"


COPY . /marketplace

WORKDIR /marketplace

RUN go mod init github.com/DjentBoiiii/marketplace && go mod tidy || go mod tidy

SHELL ["/bin/bash", "-c"] 