# Build stage with Bun
FROM oven/bun:1 AS builder

# Set working directory
WORKDIR /app

# Copy package files (Bun uses bun.lockb for locking)
COPY package.json bun.lockb ./

# Install dependencies (frozen-lockfile for reproducible build)
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application (Vite outputs to dist/)
RUN bun run build

# Production stage with Nginx
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/health || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]