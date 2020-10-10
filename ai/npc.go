package ai

import (
	"github.com/faiface/pixel"
)

type Npc struct {
	Counter float64
}

func (npc *Npc) Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, colliders *[]interface{}, playerPos *pixel.Vec, ec *pixel.Circle) (pixel.Vec, float64) {
	return entityPos, dir
}
