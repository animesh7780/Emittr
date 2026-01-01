# Backend Dockerfile
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN GONOSUMDB="*" go mod download -x

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o backend .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/backend .

ENV PORT=8080
EXPOSE 8080

CMD ["./backend"]
