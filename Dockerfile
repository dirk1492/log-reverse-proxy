FROM golang:alpine

RUN apk update && apk add gcc musl-dev upx ca-certificates dep git

WORKDIR /go/src/github.com/dirk1492/log-reverse-proxy

COPY * ./

RUN dep ensure
RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o /service
RUN upx /service

FROM scratch
COPY --from=0 /service /service
CMD ["/service"]