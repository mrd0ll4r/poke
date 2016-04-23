package poke

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestNewInfohash(t *testing.T) {
	ih := NewInfohash(rand.New(rand.NewSource(0)))
	assert.Equal(t, 20, len(ih))

	ih2 := NewInfohash(rand.New(rand.NewSource(0)))
	assert.NotEqual(t, ih, ih2)
}

func TestNewPeer(t *testing.T) {
	peer := NewPeer(rand.New(rand.NewSource(0)))
	assert.Equal(t, 20, len(peer.ID))
	assert.True(t, peer.Port >= 1024)

	peer2 := NewPeer(rand.New(rand.NewSource(0)))
	assert.Equal(t, 20, len(peer2.ID))
	assert.True(t, peer2.Port >= 1024)

	assert.NotEqual(t, peer.Port, peer2.Port)
	assert.NotEqual(t, peer.ID, peer2.ID)
	assert.NotEqual(t, peer.IP, peer2.IP)
}
