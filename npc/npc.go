package npc

import (
	"math"
	"strconv"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/ai"
	"github.com/deathowl/rpg/engine"
	"github.com/faiface/pixel"
)

type animState int

const (
	idle animState = iota
)

type NPC struct {
	Sheet      pixel.Picture
	Anims      map[string][]pixel.Rect
	Rate       float64
	Ai         ai.BaseAi
	Pos        pixel.Vec
	state      animState
	Counter    float64
	Dir        float64
	vel        pixel.Vec
	frame      pixel.Rect
	Size       float64
	SpriteSize float64
	RunSpeed   float64
	sprite     *pixel.Sprite
	Collider   *pixel.Circle
	Name       string
}

func NewNPC(eobj *tiled.Object, collider *pixel.Circle) *NPC {
	var enemyAi ai.BaseAi
	var sheet pixel.Picture
	var sheetsize float64
	var anims map[string][]pixel.Rect
	var spd float64
	var name string
	for _, prop := range eobj.Properties {
		if prop.Name == "ai" {
			enemyAi = ai.GetAi(prop.Value)
		}
		if prop.Name == "sheetsize" {
			sheetsize, _ = strconv.ParseFloat(prop.Value, 64)
		}
		if prop.Name == "spritesheet" {
			sheet, anims, _ = engine.LoadAnimationSheet("assets/"+prop.Value+".png", "assets/"+prop.Value+".csv", sheetsize)
		}
		if prop.Name == "movementspeed" {
			spd, _ = strconv.ParseFloat(prop.Value, 64)
		}
		if prop.Name == "name" {
			name = prop.Value
		}

	}
	return &NPC{Ai: enemyAi, Sheet: sheet, Anims: anims, Rate: 1.0 / 10,
		Dir: +1, Pos: pixel.V(eobj.X+8, eobj.Y+8), Size: eobj.Width, SpriteSize: sheetsize, RunSpeed: spd, Collider: collider, Name: name}
}

func (npc *NPC) Update(dt float64) {
	npc.Counter += dt
	aRate := .6
	i := int(math.Floor(npc.Counter / aRate))
	npc.frame = npc.Anims["Idle"][i%len(npc.Anims["Idle"])]
}

func (npc *NPC) Draw(t pixel.Target) {
	if npc.sprite == nil {
		npc.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	npc.sprite.Set(npc.Sheet, npc.frame)
	npc.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(-npc.Dir*(npc.Size/npc.SpriteSize), 1*(npc.Size/npc.SpriteSize))).
		Moved(npc.Pos))
}
