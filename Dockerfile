# Multi-Stage production build
FROM golang:1.17.3 as deploy

# Retrieve the binary from the previous stage
COPY ./rekor-server /usr/local/bin/rekor-server

# Set the binary as the entrypoint of the container
CMD ["rekor-server", "serve"]
