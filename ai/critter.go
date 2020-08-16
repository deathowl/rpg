package ai

import (
	"github.com/faiface/pixel"

	"github.com/deathowl/go-tiled"
)

type Critter struct {
	Counter int
}

func (critter *Critter) Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, world *tiled.Map) (pixel.Vec, float64) {
	critter.Counter++
	if critter.Counter == 100 {
		dir = (-1.0 * dir)
		critter.Counter = 0
	}

	entityPos = pixel.V(entityPos.X+(dir*speed*dt), entityPos.Y)
	return entityPos, dir
}
