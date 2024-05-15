/*
 * SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	"github.com/gardener/ext-authz-server/pkg/auth"
	auth_v3 "github.com/gardener/ext-authz-server/pkg/auth/v3"
)

const service auth.Services = `^outbound\|(1194|443)\|\|(vpn-seed-server(-[0-4])?|kube-apiserver)\..*\.svc\.cluster\.local$`

func main() {
	port := flag.Int("port", 9001, "gRPC port")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen to %d: %v", *port, err)
	}

	gs := grpc.NewServer()

	envoy_service_auth_v3.RegisterAuthorizationServer(gs, auth_v3.New(service))

	log.Printf("starting gRPC server on: %d\n", *port)

	err = gs.Serve(lis)
	if err != nil {
		log.Fatalf("Server stopped with error %v", err)
	}
}
