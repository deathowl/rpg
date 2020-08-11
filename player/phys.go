package player

import (
	"fmt"

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
		if !checkcollisions(pixel.V(ctrl.X-pp.RunSpeed*dt, ctrl.Y), pp.Colliders) {
			ctrl.X -= pp.RunSpeed * dt
		}
	case RIGHT:
		pp.vel.X = +pp.RunSpeed
		if !checkcollisions(pixel.V(ctrl.X+pp.RunSpeed*dt, ctrl.Y), pp.Colliders) {
			ctrl.X += pp.RunSpeed * dt
		}
	default:
		pp.vel.X = 0
	}
	switch d {
	case UP:
		pp.vel.Y = +pp.RunSpeed
		if !checkcollisions(pixel.V(ctrl.X, ctrl.Y+pp.RunSpeed*dt), pp.Colliders) {
			ctrl.Y += pp.RunSpeed * dt
		}
	case DOWN:
		pp.vel.Y = -pp.RunSpeed
		if !checkcollisions(pixel.V(ctrl.X, ctrl.Y-pp.RunSpeed*dt), pp.Colliders) {

			ctrl.Y -= pp.RunSpeed * dt
		}
	default:
		pp.vel.Y = 0
	}

	// apply gravity and velocity
	//pp.rect = pp.rect.Moved(pp.vel.Scaled(dt))

	// check collisions against each platform
	// gp.ground = false

	// // jump if on the ground and the player wants to jump
	// if gp.ground && ctrl.Y > 0 {
	// 	gp.vel.Y = gp.jumpSpeed
	// }
	return ctrl
}

func checkcollisions(v pixel.Vec, colliders *[]interface{}) bool {
	obcollider := pixel.C(v, 10)
	//fmt.Println(colliders)
	for _, collider := range *colliders {
		switch v := collider.(type) {
		case pixel.Circle:
			if v.Intersect(obcollider).Radius != 0 {
				fmt.Println("collided with ", v)
				return true
			}
		case pixel.Rect:
			if v.IntersectCircle(obcollider) != pixel.ZV {
				fmt.Println("collided with ", v)
				return true
			}
		case pixel.Line:
			if v.IntersectCircle(obcollider) != pixel.ZV {
				fmt.Println("collided with ", v)
				return true
			}
		}
	}

	return false
}
