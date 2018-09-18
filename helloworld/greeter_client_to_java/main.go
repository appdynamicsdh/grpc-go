/*
 *
 * C pyum += sum
        fmt.Println(sum)
        time.Sleep(300 * time.Millisecond)
    }ight 2015 gRPC authors.
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
	appd "appdynamics"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"time"
)

const (
	address     = "localhost:20000"
	defaultName = "world"
)

func main() {
	cfg := appd.Config{}

	cfg.AppName = "gRPCTest"
	cfg.TierName = "goToJavaClientTier"
	cfg.NodeName = "goToJavaClientTier1"
	cfg.Controller.Host = "localhost"
	cfg.Controller.Port = 32774
	cfg.Controller.UseSSL = false
	cfg.Controller.Account = "customer1"
	cfg.Controller.AccessKey = ""
	cfg.InitTimeoutMs = 1000 // Wait up to 1s for initialization to finish

	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}

	backendName := "GRPC Go To Java"
	backendType := "HTTP"
	backendProperties := map[string]string{
		"HOST": "localhost",
		"PORT": "20000",
	}
	resolveBackend := true
	appd.AddBackend(backendName, backendType, backendProperties, resolveBackend)

	for {
		time.Sleep(1000 * time.Millisecond)

		// start the "Checkout" transaction
		btHandle := appd.StartBT("Network Connect", "")

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

		defer cancel()

		inventoryEcHandle := appd.StartExitcall(btHandle, backendName)
		hdr := appd.GetExitcallCorrelationHeader(inventoryEcHandle)
		log.Print(hdr)
		ctx = metadata.AppendToOutgoingContext(ctx, appd.APPD_CORRELATION_HEADER_NAME, hdr)

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
		appd.EndBT(btHandle)

		// end the transaction
		//appd.EndBT(btHandle)

	}
}
