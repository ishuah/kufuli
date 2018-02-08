package registry

import (
	"log"
	"sync"
	"time"

	"github.com/ishuah/kufuli/config"
)

// Register holds all the locks
type Register struct {
	state  map[string]*Service
	config config.Config
	sync.RWMutex
}

// Service holds service identifier and expiry time
type Service struct {
	ServiceID string
	Expiry    time.Time
}

// NewRegister returns a new Register object instance
func NewRegister(config config.Config) *Register {
	return &Register{state: make(map[string]*Service), config: config}
}

func (rg *Register) load(key string) (*Service, bool) {
	rg.RLock()
	defer rg.RUnlock()
	service, loaded := rg.state[key]
	return service, loaded
}

func (rg *Register) store(key string, service *Service) {
	rg.Lock()
	defer rg.Unlock()
	rg.state[key] = service
}

func (rg *Register) loadOrStore(key, value string) bool {
	_, loaded := rg.load(key)
	if !loaded {
		service := &Service{ServiceID: value, Expiry: time.Now().Add(rg.config.MaxLockSpan)}
		rg.store(key, service)
	}
	return loaded
}

func (rg *Register) delete(key string) {
	rg.Lock()
	defer rg.Unlock()
	delete(rg.state, key)
}

// FilterStaleLocks looks for locks that have expired and passes them to CleanUpStaleLocks
func (rg *Register) FilterStaleLocks(staleLocks chan string) {
	for {
		time.Sleep(5 * time.Second)
		t := time.Now()
		rg.RLock()
		for key, value := range rg.state {
			if t.After(value.Expiry) {
				log.Printf("Cleaning up long running lock on %s\n", key)
				staleLocks <- key
			}
		}
		rg.RUnlock()
	}
}

// CleanUpStaleLocks deletes expired locks
func (rg *Register) CleanUpStaleLocks(staleLocks chan string) {
	for {
		select {
		case key := <-staleLocks:
			rg.delete(key)
		}
	}
}

// LockResource locks a resource to the specified service
// Returns true on success, false if another service has locked
// the resource
func (rg *Register) LockResource(resource, serviceID string) bool {
	loaded := rg.loadOrStore(resource, serviceID)
	return !loaded
}

// ReleaseResource releases a previously locked resource
func (rg *Register) ReleaseResource(resource string) {
	rg.delete(resource)
}
