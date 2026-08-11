FROM golang:1.25.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ecommerce-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o ecommerce-seed ./cmd/seed

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates wget

COPY --from=builder /app/ecommerce-api /app/ecommerce-api
COPY --from=builder /app/ecommerce-seed /app/ecommerce-seed
COPY --from=builder /app/docs /app/docs

RUN mkdir -p /app/uploads/products

EXPOSE 8080

ENTRYPOINT ["/app/ecommerce-api"]