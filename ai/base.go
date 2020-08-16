package ai

import (
	"github.com/deathowl/go-tiled"
	"github.com/faiface/pixel"
)

type BaseAi interface {
	Tick(dt float64, entityPos *pixel.Vec, world *tiled.Map)
}

func GetAi(aistr string) BaseAi {
	switch aistr {
	case "critter":
		return &Critter{}
	}
	panic("INVALID AI TYPE SPECIFIED IN TILEMAP")
}
