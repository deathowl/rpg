package player

import (
	"encoding/csv"
	"image"
	_ "image/png"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

func LoadAnimationSheet(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the spritesheet
	sheetFile, err := os.Open(sheetPath)
	if err != nil {
		return nil, nil, err
	}
	defer sheetFile.Close()
	sheetImg, _, err := image.Decode(sheetFile)
	if err != nil {
		return nil, nil, err
	}
	sheet = pixel.PictureDataFromImage(sheetImg)

	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+frameWidth <= sheet.Bounds().Max.X; x += frameWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+frameWidth,
			sheet.Bounds().H(),
		))
	}

	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		anims[name] = frames[start : end+1]
	}

	return sheet, anims, nil
}

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
	if phys.vel.X == 0 && phys.vel.Y < 0 {
		newState = walkup
	}
	if phys.vel.X == 0 && phys.vel.Y > 0 {
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
