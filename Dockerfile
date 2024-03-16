ARG GO_VERSION=1.21.3
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /home/nekonotes
COPY .env .
RUN git clone https://github.com/hellspawn679/itllenotes-backend.git
EXPOSE 7000

