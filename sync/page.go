package sync

import "sync"

// LastPage is a thread-safe boolean for determine the last page.
type LastPage struct {
	t bool
	m sync.RWMutex
}

// Get retrieves the value of the boolean.
func (p *LastPage) Get() bool {
	p.m.RLock()
	p.m.RUnlock()
	return p.t
}

// Set sets the value of the boolean.
func (p *LastPage) Set(value bool) {
	p.m.Lock()
	p.t = value
	p.m.Unlock()
}
