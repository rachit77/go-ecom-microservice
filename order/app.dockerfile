FROM golang:1.24 AS build
RUN apt-get update && apt-get install -y gcc g++ make ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/rachit77/go-ecom-microservice
COPY go.mod go.sum ./
COPY account account
COPY catalog catalog
COPY order order
RUN go build -o /go/bin/app ./order/cmd/order

FROM alpine:3.18
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]