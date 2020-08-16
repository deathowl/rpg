package player

import (
	"github.com/deathowl/rpg/engine"
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
	RunSpeed  float64
	vel       pixel.Vec
	Colliders *[]interface{}
}

func (pp *PlayerPhys) Update(dt float64, ctrl pixel.Vec, d Direction) pixel.Vec {
	// apply controls
	switch d {
	case LEFT:
		pp.vel.X = -pp.RunSpeed
		if !engine.CheckCollisions(pixel.V(ctrl.X-pp.RunSpeed*dt, ctrl.Y), pp.Colliders) {
			ctrl.X -= pp.RunSpeed * dt
		}
	case RIGHT:
		pp.vel.X = +pp.RunSpeed
		if !engine.CheckCollisions(pixel.V(ctrl.X+pp.RunSpeed*dt, ctrl.Y), pp.Colliders) {
			ctrl.X += pp.RunSpeed * dt
		}
	default:
		pp.vel.X = 0
	}
	switch d {
	case UP:
		pp.vel.Y = +pp.RunSpeed
		if !engine.CheckCollisions(pixel.V(ctrl.X, ctrl.Y+pp.RunSpeed*dt), pp.Colliders) {
			ctrl.Y += pp.RunSpeed * dt
		}
	case DOWN:
		pp.vel.Y = -pp.RunSpeed
		if !engine.CheckCollisions(pixel.V(ctrl.X, ctrl.Y-pp.RunSpeed*dt), pp.Colliders) {

			ctrl.Y -= pp.RunSpeed * dt
		}
	default:
		pp.vel.Y = 0
	}
	return ctrl
}
