FROM golang:1.23.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -tags netgo -ldflags '-s -w' -o api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/api /api
# UNCOMMENT IF RUN LOCALLY
# COPY .env /app/.env

EXPOSE 8081

CMD ["/api"]
