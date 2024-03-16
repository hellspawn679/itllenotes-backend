ARG GO_VERSION=1.21.3
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Copying only the necessary files for Go modules to cache dependencies properly
COPY go.mod .
COPY go.sum .
RUN go mod download -x

ARG TARGETARCH
COPY . .

RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /bin/server .

FROM alpine:latest AS final
RUN apk --no-cache add  ca-certificates git 

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

USER appuser
COPY --from=build /bin/server /bin/

# Copying .env file into the container
COPY .env .

EXPOSE 7000

WORKDIR /home/nekonotes
COPY .env .

RUN git clone https://github.com/hellspawn679/itllenotes-backend.git
EXPOSE 7000
CMD go --version
