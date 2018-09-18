Install the Go SDK.

The Go agent must be downloaded from the http://download.appdynamics.com site.
Go Agent Supported Platforms
Operating Systems
•	Any Linux distribution based on glibc 2.5+.
•	Alpine Linux uses musl_libc. See https://wiki.alpinelinux.org/wiki/Running_glibc_programs for guidance on installing our agent on Alpine Linux.
Install the AppDynamics Go SDK
To install the SDK follow these steps:
1.	Download the Go Agent SDK distribution. 
2.	Extract the Go SDK ZIP into the Go workspace. 
When finished installing the Go SDK, you are ready to instrument your Go application using the API.

Instructions for using the agent are located here
https://docs.appdynamics.com/display/PRO45/Using+the+Go+Agent+SDK

There is an example gRPC Go application that supports correlation and has been verified to work located here:

Client: https://github.com/appdynamicsdh/grpc-go/blob/master/helloworld/greeter_client/main.go
Server: https://github.com/appdynamicsdh/grpc-go/blob/master/helloworld/greeter_server/main.go


The key pieces of code you need to pay attention to for correlation to work correctly are these.

For the server, the following lines are used to retrieve the correlation header, start and end a BT.
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		hdr := md.Get(appd.APPD_CORRELATION_HEADER_NAME)[0]
		bt := appd.StartBT("Fraud Detection", hdr)
		appd.EndBT(bt)
	}
  
For the client, the following lines inject the correlation id into the Metadata header:

		hdr := appd.GetExitcallCorrelationHeader(inventoryEcHandle)
		ctx = metadata.AppendToOutgoingContext(ctx, appd.APPD_CORRELATION_HEADER_NAME, hdr)
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
