package player

import (
	"github.com/faiface/pixel"
)

type Direction int

const (
	IDLE Direction = iota
	UP
	DOWN
	LEFT
	RIGHT
)

type PlayerPhys struct {
	RunSpeed float64
	vel      pixel.Vec
	ground   bool
}

func (pp *PlayerPhys) Update(dt float64, ctrl pixel.Vec, d Direction) {
	// apply controls
	switch d {
	case LEFT:
		pp.vel.X = -pp.RunSpeed
	case RIGHT:
		pp.vel.X = +pp.RunSpeed
	default:
		pp.vel.X = 0
	}
	switch d {
	case UP:
		pp.vel.Y = -pp.RunSpeed
	case DOWN:
		pp.vel.Y = +pp.RunSpeed
	default:
		pp.vel.Y = 0
	}

	// apply gravity and velocity
	//pp.rect = pp.rect.Moved(pp.vel.Scaled(dt))

	// check collisions against each platform
	// gp.ground = false
	// if gp.vel.Y <= 0 {
	// 	for _, p := range platforms {
	// 		if gp.rect.Max.X <= p.rect.Min.X || gp.rect.Min.X >= p.rect.Max.X {
	// 			continue
	// 		}
	// 		if gp.rect.Min.Y > p.rect.Max.Y || gp.rect.Min.Y < p.rect.Max.Y+gp.vel.Y*dt {
	// 			continue
	// 		}
	// 		gp.vel.Y = 0
	// 		gp.rect = gp.rect.Moved(pixel.V(0, p.rect.Max.Y-gp.rect.Min.Y))
	// 		gp.ground = true
	// 	}
	// }

	// // jump if on the ground and the player wants to jump
	// if gp.ground && ctrl.Y > 0 {
	// 	gp.vel.Y = gp.jumpSpeed
	// }
}
