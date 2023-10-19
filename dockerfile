FROM golang:1.19

WORKDIR /app


COPY go.mod go.sum ./
COPY *.go ./
COPY cmd/api/*.go ./cmd/api/
COPY internal/cart/*.go ./internal/cart/
COPY internal/order/*.go ./internal/order/
COPY internal/product/*.go ./internal/product/
COPY internal/repository/*.go ./internal/repository/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o trafilea-cart .

EXPOSE 8080

# Run
CMD ["/app/trafilea-cart"]