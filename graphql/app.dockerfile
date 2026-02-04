FROM golang:1.24 AS build
RUN apt-get update && apt-get install -y gcc g++ make ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/github.com/rachit77/go-ecom-microservice
COPY go.mod go.sum ./
COPY account account
COPY catalog catalog
COPY order order
COPY graphql graphql
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags "-s -w" -o /go/bin/app ./graphql

FROM alpine:3.18
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]