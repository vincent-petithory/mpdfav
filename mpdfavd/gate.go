package main

type GateWaiter chan bool

type Gate struct {
	waiters []GateWaiter
	opened  bool
}

func (g Gate) Opened() bool {
	return g.opened
}

func (g *Gate) Waiter() GateWaiter {
	gw := make(GateWaiter)
	g.waiters = append(g.waiters, gw)
	return gw
}

func (g *Gate) Open() bool {
	if !g.opened {
		for _, waiter := range g.waiters {
			close(waiter)
		}
		g.opened = true
	}
	return g.opened
}

func NewGate() Gate {
	return Gate{make([]GateWaiter, 0), false}
}
