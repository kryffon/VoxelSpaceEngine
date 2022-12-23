package main

import "math"

type Player struct {
	x       float64
	y       float64
	phi     float64
	height  float64
	horizon float64
	moved   bool
}

func (p *Player) Rotate(dphi float64) {
	p.phi += dphi
	if p.phi < 0 {
		p.phi += 2*math.Pi
	}
	if p.phi > 2*math.Pi {
		p.phi -= 2*math.Pi
	}
	p.moved = true
}

func (p *Player) Move(dir float64) {
	// dir = +ve W, dir = -ve S
	p.x -= dir * math.Sin(p.phi)
	p.y -= dir * math.Cos(p.phi)
	p.moved = true
}

func (p *Player) ChangeHeight(dh float64) {
	// dh = +ve R, dh = -ve F
	p.height += dh
	p.moved = true
}

func (p *Player) ChangePitch(dp float64) {
	// dp = +ve Q, dp = -ve E
	p.horizon += dp
	p.moved = true
}
