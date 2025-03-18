# 1단계: 빌드 환경 (Go 환경)
FROM golang:1.23 AS builder
WORKDIR /app

# 모듈 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 소스 코드 복사 및 빌드
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /collector ./cmd/collector

# 2단계: 실행 환경 (최소한의 이미지)
FROM alpine:latest
WORKDIR /root/

# 실행 파일 복사
COPY --from=builder /collector .

# 실행 명령어
CMD ["./collector"]