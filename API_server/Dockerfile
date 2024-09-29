FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .
COPY ./static_files/ .
COPY ./certs .

RUN go build -o API_server main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/API_server /build/API_server
COPY --from=builder /build/dicc.txt /build/dicc.txt
COPY --from=builder /build/certs /build/certs

CMD ["./API_server"]