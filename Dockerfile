FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

RUN addgroup -S app && adduser -S app -G app

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/api/openapi.yaml ./api/

RUN chown -R app:app ./

USER app

EXPOSE 8080

CMD ["./main"]