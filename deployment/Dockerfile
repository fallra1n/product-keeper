FROM golang:1.22.0-alpine

WORKDIR /app

COPY ./ ./

RUN apk update
RUN apk add postgresql-client

RUN go mod download
RUN go build -o product-keeper ./cmd/product-keeper/main.go

CMD ["./product-keeper"]
