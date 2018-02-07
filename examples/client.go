package main

import (
	"log"
	"os"

	"github.com/ishuah/kufuli/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatal("You need to provide three arguments")
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := api.NewApiClient(conn)

	action := args[0]
	var response *api.Response

	switch action {
	case "lock":
		response, err = c.RequestLock(context.Background(), &api.Request{Resource: args[1], ServiceID: args[2]})
	case "release":
		response, err = c.ReleaseLock(context.Background(), &api.Request{Resource: args[1], ServiceID: args[2]})
	}

	if err != nil {
		log.Fatalf("Error when calling %s: %s", action, err)
	}

	log.Printf("Response from server: %v", response)
}
