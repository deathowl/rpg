package ai

import (
	"github.com/faiface/pixel"
)

type Npc struct {
	Counter float64
}

func (npc *Npc) Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, colliders *[]interface{}, playerPos *pixel.Vec, ec *pixel.Circle) (pixel.Vec, float64, pixel.Vec) {
	return entityPos, dir, pixel.V(0, 0)
}
