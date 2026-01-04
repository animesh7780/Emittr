# Root Dockerfile - doesn't build anything, just a placeholder
# Railway will use individual Dockerfiles from subdirectories
FROM alpine:latest
RUN echo "Use specific service Dockerfiles in backend/, frontend/, etc."
