package health

import (
	"testing"
	"time"
)

func Test_CachedChecker_Start_Stop(t *testing.T) {
	c := NewCachedChecker()
	c.AddChecker("upTestChecker", upTestChecker{})
	c.AddChecker("upTestChecker", upTestChecker{})
	c.AddInfo("key", "value")
	c.Start(200 * time.Millisecond)
	defer c.Stop()
	time.Sleep(200 * time.Millisecond)
	health := c.Check()
	if !health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), true)
	}
}

func Test_CachedChecker_Check_Down_combined(t *testing.T) {
	c := NewCachedChecker()
	c.AddChecker("downTestChecker", downTestChecker{})
	c.AddChecker("upTestChecker", upTestChecker{})
	c.Start(200 * time.Millisecond)
	defer c.Stop()
	health := c.Check()
	if !health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), true)
	}
}
