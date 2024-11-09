# syntax=docker/dockerfile:1
#Build stage
FROM golang:1.22.3-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags "-s -w" -o main .

#Run stage
FROM alpine
RUN apk add --no-cache tzdata openssl
ENV TZ=Asia/Bangkok

WORKDIR /app

RUN addgroup --system --gid 1001 golanggroup
RUN adduser --system --uid 1001 golang

COPY --from=builder --chown=golang:golanggroup /app/main .
COPY --chown=golang:golanggroup ./app.env .

RUN mkdir -p /app/image && \
    chown golang:golanggroup /app/image && \
    chmod 755 /app/image

USER golang

EXPOSE 8080
CMD ["/app/main"] 