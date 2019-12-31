package listeners

import (
	"fmt"

	"net"
	"sync"

	"github.com/mochi-co/mqtt/server/listeners/auth"
	"github.com/mochi-co/mqtt/server/system"
)

// MockCloser is a function signature which can be used in testing.
func MockCloser(id string) {}

// MockEstablisher is a function signature which can be used in testing.
func MockEstablisher(id string, c net.Conn, ac auth.Controller) error {
	return nil
}

// MockListener is a mock listener for establishing client connections.
type MockListener struct {
	sync.RWMutex
	id          string
	Config      *Config
	address     string
	IsListening bool
	IsServing   bool
	done        chan bool
	errListen   bool
}

// NewMockListener returns a new instance of MockListener
func NewMockListener(id, address string) *MockListener {
	return &MockListener{
		id:      id,
		address: address,
		done:    make(chan bool),
	}
}

// Serve serves the mock listener.
func (l *MockListener) Serve(establisher EstablishFunc) {
	l.Lock()
	l.IsServing = true
	l.Unlock()
	for {
		select {
		case <-l.done:
			return
		}
	}
}

// SetConfig sets the configuration values of the mock listener.
func (l *MockListener) Listen(s *system.Info) error {
	if l.errListen {
		return fmt.Errorf("listen failure")
	}

	l.Lock()
	l.IsListening = true
	l.Unlock()
	return nil
}

// SetConfig sets the configuration values of the mock listener.
func (l *MockListener) SetConfig(config *Config) {
	l.Lock()
	l.Config = config
	l.Unlock()
}

// ID returns the id of the mock listener.
func (l *MockListener) ID() string {
	l.RLock()
	id := l.id
	l.RUnlock()
	return id
}

// Close closes the mock listener.
func (l *MockListener) Close(closer CloseFunc) {
	l.Lock()
	defer l.Unlock()
	l.IsServing = false
	closer(l.id)
	close(l.done)
}