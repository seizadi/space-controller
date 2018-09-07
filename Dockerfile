FROM golang:1.10-alpine as builder

RUN apk --update add --no-cache --virtual .build-deps \
    gcc libc-dev linux-headers

ENV SRC=/go/src/github.com/seizadi/space-controller

COPY . ${SRC}
WORKDIR ${SRC}

RUN go build -o bin/space-controller .

FROM alpine:3.5

ENV SRC=/go/src/github.com/seizadi/space-controller
COPY --from=builder ${SRC}/bin/space-controller /

# Allows to verify certificates
RUN apk update --no-cache && apk add --no-cache ca-certificates

ENTRYPOINT ["/space-controller"]