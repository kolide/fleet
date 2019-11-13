FROM gcr.io/distroless/base-debian10:nonroot
LABEL author Kolide Developers <engineering@kolide.co>
USER nonroot

COPY ./build/binary-bundle/linux/fleet ./build/binary-bundle/linux/fleetctl /usr/bin/

EXPOSE 8080
CMD ["/usr/bin/fleet", "serve"]
