package enemy

import (
	"fmt"
	"math"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/ai"
	"github.com/deathowl/rpg/engine"
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

type Enemy struct {
	Sheet   pixel.Picture
	Anims   map[string][]pixel.Rect
	Rate    float64
	Ai      ai.BaseAi
	Pos     pixel.Vec
	state   animState
	Counter float64
	Dir     float64
	vel     pixel.Vec
	frame   pixel.Rect

	sprite *pixel.Sprite
}

func NewEnemy(eobj *tiled.Object) Enemy {
	var enemyAi ai.BaseAi
	var sheet pixel.Picture
	var anims map[string][]pixel.Rect
	for _, prop := range eobj.Properties {
		fmt.Println(prop.Name)
		if prop.Name == "ai" {
			enemyAi = ai.GetAi(prop.Value)
		}
		if prop.Name == "spritesheet" {
			sheet, anims, _ = engine.LoadAnimationSheet("assets/"+prop.Value+".png", "assets/"+prop.Value+".csv", 64)
		}

	}
	return Enemy{Ai: enemyAi, Sheet: sheet, Anims: anims, Rate: 1.0 / 10,
		Dir: +1, Pos: pixel.V(eobj.X+8, eobj.Y+8)}
}

func (enemy *Enemy) Update(dt float64, world *tiled.Map) {
	enemy.Ai.Tick(dt, &enemy.vel, world)
	enemy.Counter += dt

	// determine the new animation state
	var newState animState
	var aRate float64
	switch {
	case enemy.vel.Len() == 0:
		newState = idle
		aRate = .6
	case enemy.vel.Len() > 0:
		newState = walk
		aRate = enemy.Rate
	}
	if enemy.vel.X == 0 && enemy.vel.Y > 0 {
		newState = walkup
	}
	if enemy.vel.X == 0 && enemy.vel.Y < 0 {
		newState = walkdown
	}

	// reset the time counter if the state changed
	if enemy.state != newState {
		enemy.state = newState
		enemy.Counter = 0
	}

	// determine the correct animation frame
	i := int(math.Floor(enemy.Counter / aRate))
	switch enemy.state {
	case idle:
		enemy.frame = enemy.Anims["Idle"][i%len(enemy.Anims["Idle"])]
	case walk:
		enemy.frame = enemy.Anims["Walk"][i%len(enemy.Anims["Walk"])]
	case walkup:
		enemy.frame = enemy.Anims["WalkUp"][i%len(enemy.Anims["WalkUp"])]
	case walkdown:
		enemy.frame = enemy.Anims["WalkDown"][i%len(enemy.Anims["WalkDown"])]

	}

	// set the facing direction of the gopher
	if enemy.vel.X != 0 {
		if enemy.vel.X > 0 {
			enemy.Dir = +1
		} else {
			enemy.Dir = -1
		}
	}
}

func (enemy *Enemy) Draw(t pixel.Target) {
	if enemy.sprite == nil {
		enemy.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	enemy.sprite.Set(enemy.Sheet, enemy.frame)
	enemy.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(-enemy.Dir, 1)).
		Moved(enemy.Pos))
}
