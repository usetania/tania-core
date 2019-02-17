FROM golang:1.11 AS golang

RUN go get github.com/tools/godep

COPY . /go/src/github.com/Tanibox/tania-core
WORKDIR /go/src/github.com/Tanibox/tania-core
RUN go get


FROM node:8-jessie as node

COPY --from=golang /go /go

WORKDIR /go/src/github.com/Tanibox/tania-core

RUN cp conf.json.example conf.json && apt-get update && apt-get install --no-install-recommends -yf dh-autoreconf=* libtool=* nasm=* && \
    npm install && npm run prod


FROM golang:1.11 as build

COPY --from=golang /go /go

WORKDIR /go/src/github.com/Tanibox/tania-core
RUN go build


FROM golang:1.11

COPY --from=build /go/ /go

WORKDIR /go/src/github.com/Tanibox/tania-core

EXPOSE 8080
ENTRYPOINT [ "/go/src/github.com/Tanibox/tania-core/tania-core" ]

