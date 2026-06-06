FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tennisdaily-api ./cmd/api

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai
ENV PORT=8080

COPY --from=builder /tennisdaily-api /app/tennisdaily-api
COPY --from=builder /app/config.yaml /app/config.yaml

EXPOSE 8081

CMD ["/app/tennisdaily-api"]
