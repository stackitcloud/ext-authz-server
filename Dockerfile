# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.21.3 AS builder

COPY . /app
WORKDIR /app
ENV CGO_ENABLED=0
RUN go install -mod=vendor ./...

FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /

COPY --from=builder /go/bin/ext-authz-server /app/server
CMD ["/app/server"]
