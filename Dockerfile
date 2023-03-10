# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

ARG TARGETARCH
ARG GIT_HEAD_COMMIT
ARG GIT_TAG_COMMIT
ARG GIT_LAST_TAG
ARG GIT_MODIFIED
ARG GIT_REPO
ARG BUILD_DATE

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on go build \
        -gcflags "-N -l" \
        -ldflags "-X main.GitRepo=$GIT_REPO -X main.GitTag=$GIT_LAST_TAG -X main.GitCommit=$GIT_HEAD_COMMIT -X main.GitDirty=$GIT_MODIFIED -X main.BuildTime=$BUILD_DATE" \
        -o manager

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER nonroot:nonroot

ENTRYPOINT ["/manager"]

FROM golang:1.19-alpine as dlv
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "--", "/manager"

