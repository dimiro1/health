package health

import (
	"context"
	"sync"
)

// Checker is a interface used to provide an indication of application health.
type Checker interface {
	Check(ctx context.Context) Health
}

// CheckerFunc is an adapter to allow the use of
// ordinary go functions as Checkers.
type CheckerFunc func(ctx context.Context) Health

func (f CheckerFunc) Check(ctx context.Context) Health {
	return f(ctx)
}

type checkerItem struct {
	name    string
	checker Checker
}

// CompositeChecker aggregate a list of Checkers
type CompositeChecker struct {
	checkers []checkerItem
	info     map[string]interface{}
}

// NewCompositeChecker creates a new CompositeChecker
func NewCompositeChecker() CompositeChecker {
	return CompositeChecker{}
}

// AddInfo adds a info value to the Info map
func (c *CompositeChecker) AddInfo(key string, value interface{}) *CompositeChecker {
	if c.info == nil {
		c.info = make(map[string]interface{})
	}

	c.info[key] = value

	return c
}

// AddChecker add a Checker to the aggregator
func (c *CompositeChecker) AddChecker(name string, checker Checker) {
	c.checkers = append(c.checkers, checkerItem{name: name, checker: checker})
}

// Check returns the combination of all checkers added
// if some check is not up, the combined is marked as down
func (c CompositeChecker) Check(ctx context.Context) Health {
	health := NewHealth()
	health.Up()

	healths := make(map[string]interface{})

	type state struct {
		h    Health
		name string
	}
	ch := make(chan state, len(c.checkers))
	var wg sync.WaitGroup
	for _, item := range c.checkers {
		wg.Add(1)
		item := item
		go func() {
			ch <- state{h: item.checker.Check(ctx), name: item.name}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	for s := range ch {
		if !s.h.IsUp() && !health.IsDown() {
			health.Down()
		}
		healths[s.name] = s.h
	}

	health.info = healths

	// Extra Info
	for key, value := range c.info {
		health.AddInfo(key, value)
	}
	return health
}
