FROM golang:1.17-alpine3.14 AS builder

WORKDIR /src

COPY . .

RUN go get -d -v ./...
RUN go build -v -o rest-api-server services/rest-api-server/main.go

FROM alpine

COPY --from=builder /src/rest-api-server /opt/app/rest-api-server

CMD ["/opt/app/rest-api-server"]
