package ai

import (
	"fmt"

	"github.com/faiface/pixel"

	"github.com/deathowl/go-tiled"
)

type Critter struct {
}

func (critter *Critter) Tick(dt float64, entityPos *pixel.Vec, world *tiled.Map) {
	fmt.Println("AI ticked")
}
