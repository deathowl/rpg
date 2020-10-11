package ai

import (
	"github.com/faiface/pixel"
)

type BaseAi interface {
	Tick(dt float64, entityPos pixel.Vec, dir float64, speed float64, colliders *[]interface{}, playerPos *pixel.Vec, ec *pixel.Circle) (pixel.Vec, float64, pixel.Vec)
}

func GetAi(aistr string) BaseAi {
	switch aistr {
	case "critter":
		return &Critter{}
	case "npc":
		return &Npc{}
	case "square1":
		return &Square1{}
	}
	panic("INVALID AI TYPE SPECIFIED IN TILEMAP")
}
