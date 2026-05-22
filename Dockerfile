# =============================================================================
# Dockerfile — go-web-app
# =============================================================================
#
# This Dockerfile uses a two-stage build:
#
#   Stage 1 (builder) — compiles the Go source code into a binary.
#   Stage 2 (final)   — copies only the binary into a tiny runtime image.
#
# Why two stages?
#   The Go compiler and all source code are only needed at build time.
#   They should never ship to production. By using two stages, the final
#   image contains only what the app needs to run — nothing more.
#
#   Full Go image    ≈ 800 MB
#   Final image      ≈  20 MB
#
# A smaller image means:
#   - Faster downloads from the container registry
#   - A much smaller attack surface if a vulnerability is found
#   - Lower data transfer costs in cloud environments
# =============================================================================

# =============================================================================
# Stage 1: builder
# =============================================================================
# The golang alpine image includes the Go compiler and standard library.
# Alpine is a minimal Linux distribution, so the builder image is also small.
FROM golang:1.22.5-alpine3.20 AS builder

# OCI image labels — these appear when you run `docker inspect` and are
# indexed by container registries for discoverability.
LABEL org.opencontainers.image.title="go-web-app" \
      org.opencontainers.image.description="Lightweight Go web server serving a course catalogue" \
      org.opencontainers.image.source="https://github.com/DevSars24/go-web-app" \
      org.opencontainers.image.licenses="MIT"

# Set the working directory inside the container.
# All subsequent COPY, RUN, and CMD instructions run relative to this path.
WORKDIR /app

# Copy go.mod first, then download dependencies.
# Docker caches each instruction as a layer. Because go.mod changes far less
# often than application code, downloading dependencies is cached separately.
# This means `go mod download` only re-runs when go.mod actually changes,
# saving 30-60 seconds on most builds.
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code and build the binary.
COPY . .

# Build flags:
#   CGO_ENABLED=0    — disable C interop so the binary has no system library
#                      dependencies. This is required for distroless images.
#   GOOS=linux       — target Linux (in case you build on macOS or Windows)
#   GOARCH=amd64     — target x86-64 CPUs (most cloud VMs)
#   -ldflags="-w -s" — strip debug symbols to reduce binary size by ~30%
#   -o main          — name the output binary "main"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o main .

# =============================================================================
# Stage 2: final (the image that ships to production)
# =============================================================================
# distroless/static contains only:
#   - CA certificates (needed for HTTPS outbound calls)
#   - Basic system files (/etc/passwd, /etc/group)
#
# It has no shell, no package manager, and no debugging tools.
# If an attacker somehow gets into the container, they cannot do much without
# a shell. This is the most secure base image for a static Go binary.
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy only the compiled binary and the static HTML files from the builder.
# Everything else — the Go toolchain, source code, and module cache — is left
# behind in the builder stage and never makes it to this image.
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

# Run as a non-root user (uid=65532, provided by the distroless image).
# Running containers as root is a security risk. If the container were
# compromised, a root process would have much more access to the host.
USER nonroot:nonroot

# Tell Docker (and Kubernetes) that the app listens on port 8080.
# This is documentation only — it does not actually open the port.
# The Kubernetes Service resource handles that.
EXPOSE 8080

# Start the server. Using the JSON array form (exec format) ensures the binary
# receives OS signals directly — which is needed for graceful shutdown to work.
ENTRYPOINT ["./main"]