/* Copyright (C) 2013 Vincent Petithory <vincent.petithory@gmail.com>
 *
 * This file is part of mpdfav.
 *
 * mpdfav is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mpdfav is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mpdfav.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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
