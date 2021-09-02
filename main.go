/*
 * Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"

	"github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service/pkg/auth"
	auth_v3 "github.com/envoyproxy/envoy/examples/ext_authz/auth/grpc-service/pkg/auth/v3"
)

const service auth.Services = `^outbound\|1194\|*vpn-seed-server\..*\.svc\.cluster\.local$`

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

	gs.Serve(lis)
}
