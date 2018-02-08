package api

import (
	"log"
	"time"

	"github.com/ishuah/kufuli/config"
	"github.com/ishuah/kufuli/registry"
	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	register *registry.Register
	config   config.Config
}

// NewServer returns a Server instance
func NewServer() Server {
	c, err := config.GetConfig()
	if err != nil {
		log.Printf("Error: %v. Falling back to default values", err)
	}
	return Server{register: registry.NewRegister(c), config: c}
}

// attemptLockWithTimeout runs with retries in case resource is already locked
func (s *Server) attemptLockWithTimeout(r *Request) (bool, string) {
	success := false

	for i := 0; i < s.config.MaxRetries; i++ {
		if success = s.register.LockResource(r.Resource, r.ServiceID); success {
			break
		}
		log.Printf("failed to acquire lock %v, retrying in 6s", r)
		time.Sleep(s.config.RetryDelay)
	}
	if !success {
		return success, "Error: timeout while trying to acquire lock"
	}
	return success, ""
}

// RequestLock locks a resource for the specified service
func (s *Server) RequestLock(ctx context.Context, r *Request) (*Response, error) {
	log.Printf("Received lock request for %s", r.Resource)

	success, err := s.attemptLockWithTimeout(r)

	return &Response{Success: success, Error: err}, nil
}

// ReleaseLock releases the specified resource
func (s *Server) ReleaseLock(ctx context.Context, r *Request) (*Response, error) {
	log.Printf("Received a release request for %s", r.Resource)
	s.register.ReleaseResource(r.Resource)
	return &Response{Success: true}, nil
}

// BackgroundWorkers starts register clean up workers
func (s *Server) BackgroundWorkers() {
	staleLocks := make(chan string)
	go s.register.FilterStaleLocks(staleLocks)
	go s.register.CleanUpStaleLocks(staleLocks)
}
