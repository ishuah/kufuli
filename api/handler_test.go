package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestServer(t *testing.T) {
	s := NewServer()
	var ctx context.Context
	request1 := &Request{Resource: "disk0", ServiceID: "logger"}
	request2 := &Request{Resource: "disk0", ServiceID: "sync"}

	// test RequestLock
	response, err := s.RequestLock(ctx, request1)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, response.Error, "")

	response, err = s.RequestLock(ctx, request2)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, response.Error, "Error: timeout while trying to acquire lock")

	// test ReleaseLock
	response, err = s.ReleaseLock(ctx, request1)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}
