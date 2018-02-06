package api

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	state sync.Map
}

func (s *Server) LockResource(r *Request) bool {
	_, loaded := s.state.LoadOrStore(r.Resource, r.ServiceID)
	log.Printf("success: %t", !loaded)
	return !loaded
}

func (s *Server) ReleaseResource(r *Request) {
	s.state.Delete(r.Resource)
}

// RequestLock locks a resource for the specified service
func (s *Server) RequestLock(ctx context.Context, r *Request) (*Response, error) {
	log.Printf("Received lock request for %s", r.Resource)
	success := false

	// retries in case of failure
	for i := 0; i < 5; i++ {
		if success = s.LockResource(r); success {
			break
		}
		log.Printf("failed to acquire lock %v, retrying in 6s", r)
		time.Sleep(6 * time.Second)
	}
	return &Response{Success: success}, nil
}

// ReleaseLock releases the specified resource
func (s *Server) ReleaseLock(ctx context.Context, in *Request) (*Response, error) {
	log.Printf("Received a release request for %s", in.Resource)
	s.ReleaseResource(in)
	return &Response{Success: true}, nil
}
