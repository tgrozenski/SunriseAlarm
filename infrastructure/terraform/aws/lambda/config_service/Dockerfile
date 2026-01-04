# Build stage
FROM golang:1.25.3-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bootstrap ./main.go

# Final stage
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.1 /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bootstrap /var/runtime/bootstrap
ENV PORT=8080
CMD ["bootstrap"]
