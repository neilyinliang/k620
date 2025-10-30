# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# 运行阶段
FROM scratch
WORKDIR /app
COPY --from=builder /app/app .


ENV APP_ENV=production
ENV APP_PORT=80
ENV SUB_ADDRESSES=a.mojocn.com,b.mojocn.com
ENV ALLOW_USERS=a420aa94-5e8a-415d-9537-484be3774daa
ENV INTERVAL_SECOND=3600
ENV ENABLE_DATA_USAGE_METERING=true

EXPOSE 80
ENTRYPOINT ["./app", "run"]