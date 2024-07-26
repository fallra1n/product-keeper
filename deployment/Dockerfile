FROM golang:1.22.0

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o /app/product-keeper ./cmd/product-keeper/main.go

CMD ["/app/product-keeper"]
