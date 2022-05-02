# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine AS builder

RUN apk --no-cache add make && \
    apk add git
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build

FROM alpine:3.15.4
COPY --from=builder /app/ext-authz-server /app/server
CMD ["/app/server"]
