FROM golang:1.23-bookworm AS builder

WORKDIR /app
COPY . .

RUN go mod download

RUN go build -o granitex .

FROM golang:1.23-bookworm AS runner

WORKDIR /app
COPY --from=builder /app/granitex .
COPY ./client /app/client
CMD ["./granitex"]