# Runner image
FROM asia-southeast2-docker.pkg.dev/flip-prod-sre-e490/flip-main/alpine

# Add packages to set timezone to WIB instead of default UTC value
RUN apk update && \
    apk add --no-cache tzdata

# Set working directory for docker container to /app
WORKDIR /app

# Create directory for storing application logs
RUN mkdir /app/log

# Copy application binary and respective config files to /app
COPY blogapp /app/main

ENTRYPOINT /app/main start
