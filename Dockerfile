# to build this docker image:
#   docker build .
FROM gocv/opencv:4.6.0 AS builder

ENV GOPATH /go

COPY . /go/src/gocv.io/x/gocv/

WORKDIR /go/src/gocv.io/x/gocv

#ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /build/qr-reader .

FROM scratch

# Copy binary and config files from /build 
# to root folder of scratch container.
COPY --from=builder ["/build/qr-reader", "/"]

ARG RTSP
ENV RTSP = ${RTSP}

CMD ["/build/qr-reader"]