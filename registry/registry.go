package registry

import "sync"

// Register holds all the locks
type Register struct {
	state map[string]string
	sync.RWMutex
}

// NewRegister returns a new Register object instance
func NewRegister() *Register {
	return &Register{state: make(map[string]string)}
}

func (rg *Register) load(key string) (string, bool) {
	rg.RLock()
	defer rg.RUnlock()
	value, loaded := rg.state[key]
	return value, loaded
}

func (rg *Register) store(key, value string) {
	rg.Lock()
	defer rg.Unlock()
	rg.state[key] = value
}

func (rg *Register) loadOrStore(key, value string) bool {
	_, loaded := rg.load(key)
	if !loaded {
		rg.store(key, value)
	}
	return loaded
}

func (rg *Register) delete(key string) {
	rg.Lock()
	defer rg.Unlock()
	delete(rg.state, key)
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
