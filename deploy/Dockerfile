FROM golang:1.24 AS build

WORKDIR /build

COPY ./.git/ /build/.git/
COPY Makefile .
COPY ./vendor/ /build/vendor/
COPY ./pkg/ /build/pkg/
COPY go.mod .
COPY ./cmd/ /build/cmd/

RUN make prebuild build-linux

FROM alpine:3.21

USER nobody

WORKDIR /

COPY --from=build /build/dist/linux/standalone .

# COPY config/config.yaml /config.yaml

CMD ["/standalone", "-config", "/config.yaml"]
