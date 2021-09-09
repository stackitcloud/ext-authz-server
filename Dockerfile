# SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

FROM golang:alpine AS builder

RUN apk --no-cache add make
COPY . /app
WORKDIR /app
RUN go mod download
RUN  go build

FROM alpine:3.13.5
COPY --from=builder /app/grpc-service /app/server
CMD ["/app/server"]