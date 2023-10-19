package main

// import (
// 	"log"
// 	"net/rpc"

// 	r "github.com/chandiniv1/go-da/da"
// )

// func main() {
// 	// Connect to the `rollkit` RPC server.
// 	client, err := rpc.Dial("tcp", "localhost:1234")
// 	if err != nil {
// 		log.Fatal("Error connecting to RPC server:", err)
// 	}
// 	defer client.Close()

// 	// Prepare a request with multiple fields.
// 	request := r.InitRequest{}

// 	// Call the remote method SayHello on the server.
// 	err = client.Call("DataAvailability.Init", request, &response)
// 	if err != nil {
// 		log.Fatal("Error calling RPC method:", err)
// 	}

// 	// Check if there's an error in the response.
// 	if response.Error != nil {
// 		log.Printf("RPC error: %s\n", response.Error)
// 	} else {
// 		log.Printf("RPC response: %s\n", response.Message)
// 	}
// }
