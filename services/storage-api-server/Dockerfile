FROM golang:1.17-alpine3.14 AS builder

WORKDIR /src

COPY . .

RUN go get -d -v ./...
RUN go build -v -o storage-api-server services/storage-api-server/main.go

FROM alpine

COPY --from=builder /src/storage-api-server /opt/app/storage-api-server

CMD ["/opt/app/storage-api-server"]
