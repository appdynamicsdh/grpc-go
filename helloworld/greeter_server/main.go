/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	appd "appdynamics"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":20050"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		hdr := md.Get(appd.APPD_CORRELATION_HEADER_NAME)[0]
		bt := appd.StartBT("Fraud Detection", hdr)
		appd.EndBT(bt)
	}

	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {

	cfg := appd.Config{}

	cfg.AppName = "gRPCTest"
	cfg.TierName = "goServerTier"
	cfg.NodeName = "goServerTier1"
	cfg.Controller.Host = ""
	cfg.Controller.Port = 8080
	cfg.Controller.UseSSL = true
	cfg.Controller.Account = "customer1"
	cfg.Controller.AccessKey = "secret"
	cfg.InitTimeoutMs = 1000 // Wait up to 1s for initialization to finish
	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
