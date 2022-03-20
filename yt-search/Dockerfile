FROM golang:1.16.5-alpine3.14 AS builder

WORKDIR /usr/local/go/src

ADD app/ /usr/local/go/src

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/main/app.go

FROM alpine:3.14

COPY --from=builder /usr/local/go/src/app /
COPY --from=builder /usr/local/go/src/config.yml /

RUN --mount=type=secret,id=YT_APIKEY \
    --mount=type=secret,id=YT_CLIENT_ID \
    --mount=type=secret,id=YT_CLIENT_SECRET \
    --mount=type=secret,id=YT_ACCESS_TOKEN \
    --mount=type=secret,id=YT_REFRESH_TOKEN \
    --mount=type=secret,id=YT_RABBIT_PASSWORD \
    export YT_APIKEY=$(cat /run/secrets/YT_APIKEY) && \
    export YT_CLIENT_ID=$(cat /run/secrets/YT_CLIENT_ID) && \
    export YT_CLIENT_SECRET=$(cat /run/secrets/YT_CLIENT_SECRET) && \
    export YT_ACCESS_TOKEN=$(cat /run/secrets/YT_ACCESS_TOKEN) && \
    export YT_REFRESH_TOKEN=$(cat /run/secrets/YT_REFRESH_TOKEN) && \
    export YT_RABBIT_PASSWORD=$(cat /run/secrets/YT_RABBIT_PASSWORD) &&

CMD ["/app"]