# Build stage
FROM golang:1.22-alpine AS build
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api ./cmd/api

# Runtime
FROM alpine:3.20
ENV GIN_MODE=release
COPY --from=build /api /bin/api
EXPOSE 8080
ENTRYPOINT ["/bin/api"]
