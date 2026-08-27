ARG IMG_TAG=latest

# Compile the shentud binary
FROM golang:1.23-alpine AS shentud-builder
WORKDIR /src/app/
ENV PACKAGES="curl make git libc-dev bash file gcc linux-headers eudev-dev"
RUN apk add --no-cache $PACKAGES

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
