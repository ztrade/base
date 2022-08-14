package fsm

import "fmt"

type Rule struct {
	Name string
	Src  []string
	Dst  string
}

type EventDesc struct {
	Name string
	Src  string
	Dst  string
	Args interface{}
}

type Callback func(event EventDesc)

type FSM struct {
	state     string
	rules     []Rule
	callbacks map[string]Callback
}

func NewFSM(initia string, rules []Rule) *FSM {
	f := new(FSM)
	f.callbacks = make(map[string]Callback)
	f.state = initia
	f.rules = rules
	return f
}

func (f *FSM) SetCallback(state string, cb Callback) {
	f.callbacks[state] = cb
}

func (f *FSM) Event(event string, args ...interface{}) (err error) {
	var src, dst string
Out:
	for _, v := range f.rules {
		if v.Name == event {
			for _, s := range v.Src {
				if s == f.state {
					dst = v.Dst
					src = s
					break Out
				}
			}
		}
	}
	if dst == "" {
		err = fmt.Errorf("current state: %s,skip event:%s", f.state, event)
		return
	}
	f.state = dst
	cb, ok := f.callbacks[dst]
	if !ok {
		return
	}
	cb(EventDesc{Name: event, Src: src, Dst: dst, Args: args})
	return
}

func (f *FSM) Current() string {
	return f.state
}
