# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package.json ./
RUN npm install
COPY web/ .
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download 2>/dev/null || true
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

# Stage 3: Final minimal image
FROM alpine:3.21
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 tuneshift
WORKDIR /app
COPY --from=backend-builder /server .
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN chown -R tuneshift:tuneshift /app
USER tuneshift
EXPOSE 8080
CMD ["./server"]
