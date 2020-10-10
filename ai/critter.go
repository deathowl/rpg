package ai

import (
	"github.com/deathowl/rpg/engine"
	"github.com/faiface/pixel"
)

type Critter struct {
	Counter float64
}

func (critter *Critter) Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, colliders *[]interface{}, playerPos *pixel.Vec, ec *pixel.Circle) (pixel.Vec, float64) {
	critter.Counter += dt
	if critter.Counter >= 2 {
		dir = (-1.0 * dir)
		critter.Counter = 0
	}

	newcolls := make([]interface{}, 0)
	for _, c := range *colliders {
		if c != ec {
			newcolls = append(newcolls, c)
		}
	}
	playerCollider := pixel.C(*playerPos, 12)
	newcolls = append(newcolls, &playerCollider)
	if !engine.CheckCollisions(pixel.V(entityPos.X+(dir*speed*dt), entityPos.Y), &newcolls) {
		entityPos = pixel.V(entityPos.X+(dir*speed*dt), entityPos.Y)
	}
	return entityPos, dir
}
