# to build this docker image:
#   docker build .
FROM gocv/opencv:4.6.0 AS builder

ENV GOPATH /go

COPY . /go/src/gocv.io/x/gocv/

WORKDIR /go/src/gocv.io/x/gocv


ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN go build -o /build/qr-reader .

CMD ["/build/qr_reader"]