FROM golang:1.22-alpine AS BUILDER
WORKDIR /app
COPY . .
RUN go build -o bin/cachecast cmd/main.go


FROM alpine:latest AS PRODUCTION
WORKDIR /build
COPY --from=builder /app/bin/cachecast .
COPY --from=builder /app/web/index.html ./web/index.html
COPY --from=builder /app/.env.docker .
RUN export APP_ENV="docker"
EXPOSE 8080
CMD "./cachecast"

