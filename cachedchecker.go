package health

import (
	"sync"
	"time"
)

// CachedChecker is a Cached Composite Checker
// You need to Start() the Cache Updater
type CachedChecker struct {
	compositeChecker *CompositeChecker
	lastState        *lastState
	done             chan struct{}
}

type lastState struct {
	health Health
	sync.RWMutex
}

// NewCachedChecker creates a new CachedChecker
func NewCachedChecker() CachedChecker {
	c := NewCompositeChecker()
	return CachedChecker{compositeChecker: &c, done: make(chan struct{})}
}

// Start will start a background Ticker to update lastState
func (c *CachedChecker) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	lastState := lastState{health: c.compositeChecker.Check()}
	c.lastState = &lastState

	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-ticker.C:
				currentHealth := c.compositeChecker.Check()
				lastState.Lock()
				lastState.health = currentHealth
				lastState.Unlock()
			}
		}
	}()
}

// Stop the background Ticker
func (c *CachedChecker) Stop() {
	c.done <- struct{}{}
}

// AddInfo adds a info value to the Info map
func (c *CachedChecker) AddInfo(key string, value interface{}) *CachedChecker {
	c.compositeChecker = c.compositeChecker.AddInfo(key, value)
	return c
}

// AddChecker add a Checker to the aggregator
func (c *CachedChecker) AddChecker(name string, checker Checker) {
	c.compositeChecker.AddChecker(name, checker)
}

// Check returns the combination of all checkers added
// if some check is not up, the combined is marked as down
func (c CachedChecker) Check() Health {
	c.lastState.RLock()
	defer c.lastState.RUnlock()
	return c.lastState.health
}
