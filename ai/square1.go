package ai

import (
	"github.com/deathowl/rpg/engine"
	"github.com/faiface/pixel"
)

type Square1 struct {
	Counter float64
	Updown  bool
}

func (critter *Square1) Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, colliders *[]interface{}, playerPos *pixel.Vec, ec *pixel.Circle) (pixel.Vec, float64, pixel.Vec) {
	critter.Counter += dt
	if critter.Counter >= 2 {
		critter.Counter = 0
		if critter.Updown {
			critter.Updown = false
		} else {
			critter.Updown = true
			dir = (-1.0 * dir)
		}

	}

	newcolls := make([]interface{}, 0)
	for _, c := range *colliders {
		if c != ec {
			newcolls = append(newcolls, c)
		}
	}
	var newpos pixel.Vec
	var vel pixel.Vec
	if critter.Updown {
		newpos = pixel.V(entityPos.X, entityPos.Y+(dir*speed*dt))
		vel = pixel.V(0, dir*speed*dt)
	} else {
		newpos = pixel.V(entityPos.X+(dir*speed*dt), entityPos.Y)
		vel = pixel.V(dir*speed*dt, 0)

	}
	playerCollider := pixel.C(*playerPos, 12)
	newcolls = append(newcolls, &playerCollider)
	if !engine.CheckCollisions(newpos, &newcolls) {
		entityPos = newpos
	}
	return entityPos, dir, vel
}
