# Multi-stage build for Go + Svelte application using Makefile

# Stage 1: Build stage
FROM node:18-alpine

# Install specific Go version and golang-migrate
ARG GO_VERSION=1.24.5
ARG MIGRATE_VERSION=v4.16.2
RUN apk add --no-cache git make curl tar && \
    curl -L "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz" | tar -C /usr/local -xz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt && \
    curl -L "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" | tar -C /tmp -xz && \
    mv /tmp/migrate /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate


WORKDIR /app

# Verify Go installation
RUN go version

# Set Go environment variables
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# Copy package files first for better caching
COPY ui/package*.json ./ui/

# Install frontend dependencies
RUN cd ui && npm ci

# Copy all source files
COPY . .

COPY .env.docker ./ui/.env


# Build using Makefile (this will build both frontend and Go app)
RUN make build

# Expose port
EXPOSE 4000

# Command to run the application
CMD ./bin/web -env production -dsn ${GLOBECHAT_DB_DSN} -gclientid ${PUBLIC_GOOGLE_CLIENT_ID} -mediadir ./media