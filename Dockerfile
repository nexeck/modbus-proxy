FROM gcr.io/distroless/static:latest

USER nonroot

# Copy our static executable.
COPY modbus-proxy /home/nonroot/service

VOLUME /home/nonroot/config.yaml

WORKDIR /home/nonroot
ENTRYPOINT ["/home/nonroot/service"]
