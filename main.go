package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ishuah/kufuli/api"
	"google.golang.org/grpc"
)

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port 7777
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := api.Server{}
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	api.RegisterApiServer(grpcServer, &s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// import (
// 	"fmt"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// type Registry struct {
// 	state sync.Map
// }

// func (r *Registry) Assign(resource string, client int) (int, bool) {
// 	occupant, fresh := r.state.LoadOrStore(resource, client)
// 	return occupant.(int), fresh
// }

// func (r *Registry) Release(resource string) bool {
// 	r.state.Delete(resource)
// 	_, ok := r.state.Load(resource)
// 	return !ok
// }

// type Client struct {
// 	ID                 int
// 	registry           *Registry
// 	availableResources []string
// }

// func (c *Client) resource() string {
// 	n := rand.Intn(len(c.availableResources))
// 	return c.availableResources[n]
// }

// func (c *Client) run() {
// 	for {
// 		resource := c.resource()
// 		occupant, fresh := c.registry.Assign(resource, c.ID)
// 		if !fresh {
// 			fmt.Printf("client %d locked resource %s\n", c.ID, resource)
// 		} else {
// 			fmt.Printf("cliend %d failed to lock resource %s, %d blocking\n", c.ID, resource, occupant)
// 		}
// 		time.Sleep(5 * time.Second)
// 	}
// }

// func main() {

// 	registry := new(Registry)
// 	clients := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// 	resources := []string{"A", "B", "C", "D"}
// 	rand.Seed(time.Now().UnixNano())

// 	for _, id := range clients {
// 		c := Client{registry: registry, ID: id, availableResources: resources}
// 		go c.run()
// 	}

// 	for {
// 		time.Sleep(10 * time.Second)
// 		fmt.Println("Releasing resources")
// 		for _, resource := range resources {
// 			registry.Release(resource)
// 		}
// 	}
// }
