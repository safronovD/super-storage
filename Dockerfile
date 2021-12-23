# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

COPY  ./build/storage  /storage

USER nonroot:nonroot

ENTRYPOINT ["/storage"]