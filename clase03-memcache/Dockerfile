# build
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/api ./cmd/api

# run
FROM alpine:3.20
ENV ADDR=":8080" \
    MEMCACHED_ADDR="memcached:11211" \
    CACHE_TTL_SECONDS="60" \
    MONGO_URI="mongodb://mongo:27017" \
    MONGO_DB="demo" \
    MONGO_COLLECTION="items"
USER 65532:65532
EXPOSE 8080
COPY --from=build /bin/api /bin/api
ENTRYPOINT ["/bin/api"]
