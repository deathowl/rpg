package player

import (
	_ "image/png"
	"math"

	"github.com/faiface/pixel"
)

type animState int

const (
	idle animState = iota
	walk
	walkup
	walkdown
	attack
)

type PlayerAnim struct {
	Sheet pixel.Picture
	Anims map[string][]pixel.Rect
	Rate  float64

	state   animState
	Counter float64
	Dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
}

func (pa *PlayerAnim) Update(dt float64, phys *PlayerPhys) {
	pa.Counter += dt

	// determine the new animation state
	var newState animState
	var aRate float64
	switch {
	case phys.vel.Len() == 0:
		newState = idle
		aRate = .6
	case phys.vel.Len() > 0:
		newState = walk
		aRate = pa.Rate
	}
	if phys.vel.X == 0 && phys.vel.Y > 0 {
		newState = walkup
	}
	if phys.vel.X == 0 && phys.vel.Y < 0 {
		newState = walkdown
	}

	// reset the time counter if the state changed
	if pa.state != newState {
		pa.state = newState
		pa.Counter = 0
	}

	// determine the correct animation frame
	i := int(math.Floor(pa.Counter / aRate))
	switch pa.state {
	case idle:
		pa.frame = pa.Anims["LeftRight"][i%len(pa.Anims["LeftRight"])]
	case walk:
		pa.frame = pa.Anims["Walk"][i%len(pa.Anims["Walk"])]
	case walkup:
		pa.frame = pa.Anims["WalkUp"][i%len(pa.Anims["WalkUp"])]
	case walkdown:
		pa.frame = pa.Anims["WalkDown"][i%len(pa.Anims["WalkDown"])]

	}

	// set the facing direction of the gopher
	if phys.vel.X != 0 {
		if phys.vel.X > 0 {
			pa.Dir = +1
		} else {
			pa.Dir = -1
		}
	}
}

func (pa *PlayerAnim) Draw(t pixel.Target, camPos *pixel.Vec) {
	if pa.sprite == nil {
		pa.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	pa.sprite.Set(pa.Sheet, pa.frame)
	pa.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(-pa.Dir, 1)).
		Moved(*camPos),
	)
}
