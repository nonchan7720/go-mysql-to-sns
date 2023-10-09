FROM golang:1.21.1-bullseye as builder

ARG TARGETPLATFORM

ENV GO111MODULE=on \
  GOPATH=/go \
  GOBIN=/go/bin

WORKDIR /workspace

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 go build -o go-mysql-to-sns \
  && chmod +x /workspace/go-mysql-to-sns

FROM gcr.io/distroless/static:nonroot
COPY --from=builder --chown=nonroot:nonroot /workspace/go-mysql-to-sns /usr/local/bin/go-mysql-to-sns
COPY --chown=nonroot:nonroot config.yaml config.yaml
ENV TZ=Asia/Tokyo
USER 65532:65532

CMD ["go-mysql-to-sns"]
