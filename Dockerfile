# Stage 1: Build SvelteKit Frontend
FROM node:22-alpine AS frontend-builder
WORKDIR /app
# Copy only package files first for caching
COPY web/package.json web/package-lock.json* ./web/
RUN cd web && npm install
# Copy the rest of the web directory
COPY web/ ./web/
RUN cd web && npm run build

# Stage 2: Build Go Backend
FROM golang:1.23-alpine AS backend-builder
# Install build tools for CGO (required by sqlite driver)
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# We don't strictly need the source 'web' dir here, just the final build in the final stage,
# but building the Go binary doesn't require the static files at compile time 
# since they are served from the filesystem at runtime.
RUN go build -o brayat-server ./cmd/server/main.go

# Stage 3: Final Slim Image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Copy the binary
COPY --from=backend-builder /app/brayat-server .

# Copy the built frontend static files
# The Go server expects them at ./web/build
COPY --from=frontend-builder /app/web/build ./web/build

# Create data directory for volume mount (Fly.io)
RUN mkdir -p /data/photos

# Environment defaults
ENV PORT=8080
ENV DATABASE_PATH=/data/brayat.db
ENV PHOTOS_DIR=/data/photos
ENV LOG_LEVEL=info
ENV GIN_MODE=release

EXPOSE 8080
CMD ["./brayat-server"]
