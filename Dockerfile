ARG IMG_TAG=latest

# Compile the shentud binary
FROM golang:1.23-alpine AS shentud-builder
WORKDIR /src/app/
ENV PACKAGES="curl make git libc-dev bash file gcc linux-headers eudev-dev"
RUN apk add --no-cache $PACKAGES

# See https://github.com/CosmWasm/wasmvm/releases
ARG WASMVM_VERSION=v2.2.1
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep ba6cb5db6b14a265c8556326c045880908db9b1d2ffb5d4aa9f09ac09b24cecc
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep b3bd755efac0ff39c01b59b8110f961c48aa3eb93588071d7a628270cc1f2326
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN LEDGER_ENABLED=false LINK_STATICALLY=true BUILD_TAGS=muslc make build

# Add to a distroless container
FROM gcr.io/distroless/cc:$IMG_TAG
ARG IMG_TAG
COPY --from=shentud-builder /src/app/build/shentud /usr/local/bin/
EXPOSE 26656 26657 1317 9090

ENTRYPOINT ["shentud", "start"]
