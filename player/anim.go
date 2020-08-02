package player
import (
	"encoding/csv"
	"os"
	"github.com/pkg/errors"
	"github.com/faiface/pixel"
	"image"
	"io"
	"strconv"
	_ "image/png"
	"math"
)

func loadAnimationSheet(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
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
	running
	jumping
	using
	hitting
)

type playerAnim struct {
	sheet pixel.Picture
	anims map[string][]pixel.Rect
	rate  float64

	state   animState
	counter float64
	dir     float64

	frame pixel.Rect

	sprite *pixel.Sprite
}
func (pa *playerAnim) update(dt float64, phys *playerPhys) {
	pa.counter += dt

	// determine the new animation state
	var newState animState
	switch {
	case !phys.ground:
		newState = jumping
	case phys.vel.Len() == 0:
		newState = idle
	case phys.vel.Len() > 0:
		newState = running
	}

	// reset the time counter if the state changed
	if pa.state != newState {
		pa.state = newState
		pa.counter = 0
	}

	// determine the correct animation frame
	switch pa.state {
	case idle:
		pa.frame = pa.anims["Front"][0]
	case running:
		i := int(math.Floor(pa.counter / pa.rate))
		pa.frame = pa.anims["Run"][i%len(pa.anims["Run"])]
	case jumping:
		speed := phys.vel.Y
		i := int((-speed/phys.jumpSpeed + 1) / 2 * float64(len(pa.anims["Jump"])))
		if i < 0 {
			i = 0
		}
		if i >= len(pa.anims["Jump"]) {
			i = len(pa.anims["Jump"]) - 1
		}
		pa.frame = pa.anims["Jump"][i]
	}

	// set the facing direction of the gopher
	if phys.vel.X != 0 {
		if phys.vel.X > 0 {
			pa.dir = +1
		} else {
			pa.dir = -1
		}
	}
}

func (pa *playerAnim) Draw(t pixel.Target, phys *playerPhys) {
	if pa.sprite == nil {
		pa.sprite = pixel.NewSprite(nil, pixel.Rect{})
	}
	// draw the correct frame with the correct position and direction
	pa.sprite.Set(pa.sheet, pa.frame)
	pa.sprite.Draw(t, pixel.IM.
		ScaledXY(pixel.ZV, pixel.V(
			phys.rect.W()/pa.sprite.Frame().W(),
			phys.rect.H()/pa.sprite.Frame().H(),
		)).
		ScaledXY(pixel.ZV, pixel.V(-pa.dir, 1)).
		Moved(phys.rect.Center()),
	)
}