# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.21.0-bullseye
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.7.1

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-10-02
LABEL Name="senzing/go-queueing-builder" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/go-queueing

# Build go program.

WORKDIR ${GOPATH}/src/go-queueing
RUN make build

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/go-queueing/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2023-08-01
LABEL Name="senzing/go-queueing" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Copy files from prior stage.

COPY --from=go_builder "/output/linux-amd64/go-queueing" "/app/go-queueing"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/go-queueing"]
