package api

import (
	"sync"
)

// Register holds all the locks
type Register struct {
	state sync.Map
}

// LockResource locks a resource to the specified service
// Returns true on success, false if another service has locked
// the resource
func (rg *Register) LockResource(r *Request) bool {
	_, loaded := rg.state.LoadOrStore(r.Resource, r.ServiceID)
	return !loaded
}

// ReleaseResource releases a previously locked resource
func (rg *Register) ReleaseResource(r *Request) {
	rg.state.Delete(r.Resource)
}
