FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .
COPY ./certs .

RUN go get -d -v
RUN go build -o proxy_server main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/proxy_server /build/proxy_server
COPY --from=builder /build/certs /build/certs

EXPOSE 8080

CMD ["./proxy_server"]