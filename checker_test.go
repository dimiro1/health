package health

import (
	"context"
	"testing"
)

type outOfServiceTestChecker struct{}

func (c outOfServiceTestChecker) Check(ctx context.Context) Health {
	health := NewHealth()
	health.OutOfService()

	return health
}

type upTestChecker struct{}

func (c upTestChecker) Check(ctx context.Context) Health {
	health := NewHealth()
	health.Up()

	return health
}

type downTestChecker struct{}

func (c downTestChecker) Check(ctx context.Context) Health {
	health := NewHealth()
	health.Down()

	return health
}

func Test_CompositeChecker_AddChecker(t *testing.T) {
	c := NewCompositeChecker()

	if len(c.checkers) != 0 {
		t.Errorf("len(c.checkers) == %d, wants %d", len(c.checkers), 0)
	}

	c.AddChecker("testChecker", upTestChecker{})

	if len(c.checkers) != 1 {
		t.Errorf("len(c.checkers) == %d, wants %d", len(c.checkers), 1)
	}
}

func Test_CompositeChecker_Check_Up(t *testing.T) {
	c := NewCompositeChecker()
	c.AddChecker("upTestChecker", upTestChecker{})

	health := c.Check(context.Background())

	if !health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), true)
	}
}

func Test_CompositeChecker_Check_Down(t *testing.T) {
	c := NewCompositeChecker()
	c.AddChecker("downTestChecker", downTestChecker{})

	health := c.Check(context.Background())

	if !health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), true)
	}
}

func Test_CompositeChecker_Check_OutOfService(t *testing.T) {
	c := NewCompositeChecker()
	c.AddChecker("outOfServiceTestChecker", outOfServiceTestChecker{})

	health := c.Check(context.Background())

	if !health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), true)
	}
}

func Test_CompositeChecker_Check_Down_combined(t *testing.T) {
	c := NewCompositeChecker()
	c.AddChecker("downTestChecker", downTestChecker{})
	c.AddChecker("upTestChecker", upTestChecker{})

	health := c.Check(context.Background())

	if !health.IsDown() {
		t.Errorf("health.IsDown() == %t, wants %t", health.IsDown(), true)
	}
}

func Test_CompositeChecker_Check_Up_combined(t *testing.T) {
	c := NewCompositeChecker()
	c.AddChecker("upTestChecker", upTestChecker{})
	c.AddChecker("upTestChecker", upTestChecker{})

	health := c.Check(context.Background())

	if !health.IsUp() {
		t.Errorf("health.IsUp() == %t, wants %t", health.IsUp(), true)
	}
}

func Test_CompositeChecker_AddInfo(t *testing.T) {
	c := NewCompositeChecker()

	c.AddInfo("key", "value")

	_, ok := c.info["key"]

	if !ok {
		t.Error("c.AddInfo() should add a key value to the map")
	}
}

func Test_CompositeChecker_AddInfo_null_map(t *testing.T) {
	c := CompositeChecker{}

	c.AddInfo("key", "value")

	_, ok := c.info["key"]

	if !ok {
		t.Error("c.AddInfo() should add a key value to the map")
	}
}

func TestCheckerFunc_Check(t *testing.T) {
	f := CheckerFunc(func(ctx context.Context) Health {
		h := NewHealth()
		h.Up()

		return h
	})

	h := f.Check(context.Background())

	if !h.IsUp() {
		t.Errorf("h.IsUp() == %t, wants %t", h.IsUp(), true)
	}
}
