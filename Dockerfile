# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o k620

# 运行阶段
FROM scratch
WORKDIR /app
COPY --from=builder /app/k620 .

EXPOSE 8226
ENTRYPOINT ["./k620"]