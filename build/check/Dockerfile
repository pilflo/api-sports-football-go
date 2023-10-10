ARG GO_VERSION
ARG ALPINE_VERSION=3.18
ARG GO_IMAGE=golang:${GO_VERSION}-alpine${ALPINE_VERSION}

FROM ${GO_IMAGE} AS base

RUN apk update && apk add --no-cache curl git zip unzip make bash jq yq

RUN git config --global --add safe.directory '*'

# Go tools
ARG GOSEC_VERSION=v2.16.0
ARG GOVULN_VERSION=v0.1.0
ARG GOIMPORTS_VERSION=v0.9.1
ARG GOFUMPT_VERSION=v0.5.0
ARG GOCILINT_VERSION=v1.54.2

# gosec
RUN curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOSEC_VERSION}
# govulncheck
RUN go install golang.org/x/vuln/cmd/govulncheck@${GOVULN_VERSION}
# goimports
RUN go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION}
# gofumpt
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION}
# golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOCILINT_VERSION}