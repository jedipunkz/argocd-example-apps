FROM golang:1.14.2-alpine3.11
MAINTAINER jedipunkz

WORKDIR /go/src/

ADD . /go/src/

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/get-servers-num

ENTRYPOINT ["/go/bin/get-servers-num"]
