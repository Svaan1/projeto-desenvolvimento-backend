FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN go build -o main cmd/server/main.go

FROM golang:1.23-alpine AS runner

# install Air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# bins for Air
COPY --from=builder /app/main /app/tmp/main
COPY . .

CMD ["air", "-c", "configs/.docker-air.toml"]