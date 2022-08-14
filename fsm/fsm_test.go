package fsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSM(t *testing.T) {
	rules := []Rule{
		{"openOrder", []string{"init"}, "open"},
		{"stopOrder", []string{"open", "openMore"}, "init"},
		{"closeOrder", []string{"open", "openMore"}, "init"},
		{"openOrder", []string{"open"}, "openMore"},
	}
	args := []interface{}{"Hello", 1, 2, 3}
	cb := func(event EventDesc) {
		assert.Equal(t, event, EventDesc{Name: "openOrder", Src: "init", Dst: "open", Args: args}, "state error")
		t.Log(event.Name, event.Src, event.Dst, event.Args)
	}
	fsm := NewFSM("init", rules)
	fsm.SetCallback("open", cb)
	assert.Equal(t, fsm.Current(), "init", "init state error")
	fsm.Event("openOrder", args...)
	assert.Equal(t, fsm.Current(), "open", "open state error")
	err := fsm.Event("test")
	assert.NotNil(t, err, "error check failed")
	fsm.Event("closeOrder")
	assert.Equal(t, fsm.Current(), "init", "close state error")
}
