package engine

import (
	"encoding/csv"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/deathowl/go-tiled"
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

func CheckCollisions(v pixel.Vec, colliders *[]interface{}) bool {
	obcollider := pixel.C(v, 10)
	//fmt.Println(colliders)
	for _, collider := range *colliders {
		switch v := collider.(type) {
		case *pixel.Circle:
			if v.Intersect(obcollider).Radius != 0 {
				fmt.Println("collided with ", v)
				return true
			}
		case *pixel.Rect:
			if v.IntersectCircle(obcollider) != pixel.ZV {
				fmt.Println("collided with ", v)
				return true
			}
		case *pixel.Line:
			if v.IntersectCircle(obcollider) != pixel.ZV {
				fmt.Println("collided with ", v)
				return true
			}
		}
	}

	return false
}

//Fixes regarding orthogonal tiled crap
func FlipY(p *tiled.Object, totalHeight float64) {
	p.Y = totalHeight - p.Y - p.Height
}
func ScaleX(p *tiled.Object, scalingFac float64) {
	p.X = p.X * scalingFac
}
func ScaleY(p *tiled.Object, scalingFac float64) {
	p.Y = p.Y * scalingFac
}
func LFlipY(o *pixel.Line, totalHeight float64) {
	o.A.Y = totalHeight - o.A.Y - 1
	o.B.Y = totalHeight - o.B.Y - 1

}
func LScaleX(o *pixel.Line, scalingFac float64) {
	o.A.X = o.A.X * scalingFac
	o.B.X = o.B.X * scalingFac

}
func LScaleY(o *pixel.Line, scalingFac float64) {
	o.A.Y = o.A.Y * scalingFac
	o.B.Y = o.B.Y * scalingFac
}
