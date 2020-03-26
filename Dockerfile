FROM gcr.io/distroless/static:nonroot
LABEL maintainer="engineering@kolide.co"
USER nonroot

COPY ./build/binary-bundle/linux/fleet ./build/binary-bundle/linux/fleetctl /usr/bin/

EXPOSE 8080
CMD ["/usr/bin/fleet", "serve"]
