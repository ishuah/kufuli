package registry

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/ishuah/kufuli/config"
)

// Register the state map, a
type Register struct {
	state  map[string]interface{}
	config config.Config
	sync.RWMutex
}

// Service holds service identifier and expiry time
type Service struct {
	ServiceID string
	Expiry    time.Time
}

type ResourceSet struct {
	resources map[string]bool
	sync.RWMutex
}

func (rs *ResourceSet) Add(resource string) *ResourceSet {
	rs.Lock()
	defer rs.Unlock()
	if rs.resources == nil {
		rs.resources = make(map[string]bool)
	}

	_, loaded := rs.resources[resource]
	if !loaded {
		rs.resources[resource] = true
	}
	return rs
}

func (rs *ResourceSet) Delete(resource string) {
	rs.Lock()
	rs.Unlock()
	delete(rs.resources, resource)
}

func (rs *ResourceSet) Has(resource string) bool {
	rs.RLock()
	rs.RUnlock()
	_, loaded := rs.resources[resource]
	return loaded
}

// NewRegister returns a new Register object instance
func NewRegister(config config.Config) *Register {
	return &Register{state: make(map[string]interface{}), config: config}
}

func (rg *Register) load(prefix, key string) (interface{}, bool) {
	rg.RLock()
	defer rg.RUnlock()
	value, loaded := rg.state[prefix+":"+key]
	return value, loaded
}

func (rg *Register) store(prefix, key string, value interface{}) {
	rg.Lock()
	defer rg.Unlock()
	rg.state[prefix+":"+key] = value
}

func (rg *Register) loadOrStore(key, value string) bool {
	_, loaded := rg.load("resource", key)
	if !loaded {
		service := &Service{ServiceID: value, Expiry: time.Now().Add(rg.config.MaxLockSpan)}
		rg.store("resource", key, service)
		_, ok := rg.load("service", value)
		if !ok {
			rg.store("service", value, &ResourceSet{})
		}
		i, _ := rg.load("service", value)
		rs, ok := i.(*ResourceSet)
		rs.Add(key)
	}
	return loaded
}

func (rg *Register) masterDelete(key string) {
	if strings.HasPrefix(key, "resource:") {
		key = strings.Replace(key, "resource:", "", 1)
	}

	s, ok := rg.load("resource", key)

	if !ok {
		return
	}
	service := s.(*Service)
	r, ok := rg.load("service", service.ServiceID)
	rs := r.(*ResourceSet)
	rs.Delete(key)
	rg.delete("resource", key)
}

func (rg *Register) delete(prefix, key string) {
	rg.Lock()
	defer rg.Unlock()
	delete(rg.state, prefix+":"+key)
}

// FilterStaleLocks looks for locks that have expired and passes them to CleanUpStaleLocks
func (rg *Register) FilterStaleLocks(staleLocks chan string) {
	for {
		time.Sleep(5 * time.Second)
		t := time.Now()
		rg.RLock()
		for key, value := range rg.state {
			service, ok := value.(*Service)
			if !ok {
				continue
			}

			if t.After(service.Expiry) {
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
			rg.masterDelete(key)
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
	rg.masterDelete(resource)
}
