package registry

import (
	"testing"

	"github.com/ishuah/kufuli/config"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	c, _ := config.GetConfig()
	rg := NewRegister(c)

	resource0 := "disk0"
	resource1 := "disk1"
	service0 := "sync"
	service1 := "logger"

	// test LockResource
	loaded := rg.LockResource(resource0, service0)
	assert.True(t, loaded)

	loaded = rg.LockResource(resource0, service1)
	assert.False(t, loaded)

	// test ReleaseResource
	ok := rg.ReleaseResource(resource0)
	_, loaded = rg.state[resource0]
	assert.False(t, loaded)
	assert.True(t, ok)

	ok = rg.ReleaseResource(resource1)
	assert.False(t, ok)
}
