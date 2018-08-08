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

package main

import (
	"log"
	"os"
	"time"
	appd "appdynamics"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"fmt"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:20050"
	defaultName = "world"

)

func main() {
	cfg := appd.Config{}

	cfg.AppName = "gRPCTest"
	cfg.TierName = "goClientTier"
	cfg.NodeName = "goClientTier1"
	cfg.Controller.Host = ""
	cfg.Controller.Port = 8080
	cfg.Controller.UseSSL = true
	cfg.Controller.Account = "customer1"
	cfg.Controller.AccessKey = "secret"
	cfg.InitTimeoutMs = 1000  // Wait up to 1s for initialization to finish

	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}


	// start the "Checkout" transaction
	btHandle := appd.StartBT("Checkout", "")

	inventoryEcHandle := appd.StartExitcall(btHandle, "Inventory DB")
	hdr := appd.GetExitcallCorrelationHeader(inventoryEcHandle)


	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	ctx = metadata.AppendToOutgoingContext(ctx, appd.APPD_CORRELATION_HEADER_NAME, hdr)

	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)


	// set the exit call details
	if err := appd.SetExitcallDetails(inventoryEcHandle, "Exitcall Detail String"); err != nil {
		log.Print(err)
	}

	appd.EndExitcall(inventoryEcHandle)

	// end the transaction
	appd.EndBT(btHandle)
}
