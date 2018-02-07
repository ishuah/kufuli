package registry

import "sync"

// Register holds all the locks
type Register struct {
	resourceMap map[string]string
	sync.RWMutex
}

// NewRegister returns a new Register object instance
func NewRegister() *Register {
	return &Register{resourceMap: make(map[string]string)}
}

func (rg *Register) load(key string) (string, bool) {
	rg.RLock()
	value, loaded := rg.resourceMap[key]
	rg.RUnlock()
	return value, loaded
}

func (rg *Register) store(key, value string) {
	rg.Lock()
	rg.resourceMap[key] = value
	rg.Unlock()
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
	delete(rg.resourceMap, key)
	rg.Unlock()
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
