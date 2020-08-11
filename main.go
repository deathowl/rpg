package main

import (
	"fmt"
	"math"
	"time"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/player"
	"github.com/deathowl/rpg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var clearColor = colornames.Black

var (
	frames = 0
	second = time.Tick(time.Second)
)

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, initialPos *pixel.Vec, colliders *[]interface{}) {
	batches := make([]*pixel.Batch, 0)
	var (
		camPos       = *initialPos
		camSpeed     = 50.0
		camZoom      = 4.0
		camZoomSpeed = 1.2
	)

	last := time.Now()
	rsprirte := pixel.NewSprite(renderedBg, renderedBg.Bounds())
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(win.Bounds().Size().X/renderedBg.Bounds().Size().X, win.Bounds().Size().Y/renderedBg.Bounds().Size().Y))
	sheet, anims, err := player.LoadAnimationSheet("assets/sheet.png", "assets/spritesheet.csv", 12)
	panicIfErr(err)
	phys := &player.PlayerPhys{
		RunSpeed:  camSpeed,
		Colliders: colliders,
	}
	anim := &player.PlayerAnim{
		Sheet: sheet,
		Anims: anims,
		Rate:  1.0 / 10,
		Dir:   +1,
	}
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		var curdir player.Direction
		if win.Pressed(pixelgl.KeyLeft) {
			curdir = player.LEFT
		}
		if win.Pressed(pixelgl.KeyRight) {
			curdir = player.RIGHT
		}
		if win.Pressed(pixelgl.KeyDown) {
			curdir = player.DOWN
		}
		if win.Pressed(pixelgl.KeyUp) {
			curdir = player.UP
		}
		camPos = phys.Update(dt, camPos, curdir)
		if win.Pressed(pixelgl.KeySpace) {
			fmt.Println(camPos.X)
			fmt.Println(camPos.Y)
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(clearColor)

		// Draw tiles
		for _, batch := range batches {
			batch.Clear()
		}
		rsprirte.Draw(win, mat)
		// imd := imdraw.New(nil)
		// imd.SetMatrix(mat)
		// for _, b := range *boundaries {
		// 	imd.Push(b.A)
		// 	imd.Push(b.B)
		// 	imd.Line(1)
		// }
		charcolliderd := imdraw.New(nil)
		charcolliderd.Push(camPos)
		charcolliderd.Circle(8, 1)
		colliderd := imdraw.New(nil)
		colliderd.SetMatrix(mat)
		for _, collider := range *colliders {
			switch v := collider.(type) {
			case pixel.Circle:
				charcolliderd.Push(v.Center)
				charcolliderd.Circle(v.Radius, 1)
			case pixel.Rect:
				charcolliderd.Push(v.Min, v.Max)
				charcolliderd.Rectangle(1)
			case pixel.Line:
				charcolliderd.Push(v.A, v.B)
				charcolliderd.Line(1)
			}
		}
		//imd.Draw(win)
		if win.Pressed(pixelgl.KeyRightControl) {
			colliderd.Draw(win)
			charcolliderd.Draw(win)
		}

		anim.Update(dt, phys)
		anim.Draw(win, &camPos)

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("RPG | FPS: %d", frames))
			fmt.Println("RPG | FPS: ", frames)
			frames = 0
		default:
		}
		win.Update()
	}
}

func initialize() {
	fmt.Println("initialize called")
	// Create the window with OpenGL
	cfg := pixelgl.WindowConfig{
		Title:  "Tiled Rpg",
		Bounds: pixel.R(0, 0, 1280, 1024),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	panicIfErr(err)

	// Initialize art assets (i.e. the tilemap)
	tilemap := world.LoadTileMap("./assets/island.tmx")
	renderedBg := world.RenderTilemap(&tilemap)

	startPos := win.Bounds().Center()
	scalingFacX := win.Bounds().Size().X / renderedBg.Bounds().Size().X
	scalingFacY := win.Bounds().Size().Y / renderedBg.Bounds().Size().Y
	colliders := make([]interface{}, 0)
	for _, ob := range tilemap.ObjectGroups[0].Objects {
		flipY(ob, renderedBg.Bounds().Size().Y)
		scaleX(ob, scalingFacX)
		scaleY(ob, scalingFacY)
		if ob.Type == "start" {
			//startPos = pixel.V(ob.X, math.Abs(ob.Y)).ScaledXY(scalingVec)
			startPos = pixel.Vec{X: ob.X, Y: ob.Y}
		}
		if ob.Type == "border" {
			var prevPoint *tiled.Point
			points := *ob.Polygons[0].Points
			for idx, p := range points {
				PFlipY(p, renderedBg.Bounds().Size().Y)
				PScaleX(p, scalingFacX)
				PScaleY(p, scalingFacY)
				fmt.Println(p)
				if idx == 0 {
					prevPoint = p
				} else if idx == len(points)-1 {
					l1 := pixel.L(pixel.V(ob.X+prevPoint.X, ob.Y+prevPoint.Y), pixel.V(ob.X+p.X, ob.Y+p.Y))
					l2 := pixel.L(pixel.V(ob.X+p.X, ob.Y+p.Y), pixel.V(ob.X+points[0].X, ob.Y+points[0].Y))
					colliders = append(colliders, l1, l2)

				} else {
					colliders = append(colliders, pixel.L(pixel.V(ob.X+prevPoint.X, ob.Y+prevPoint.Y), pixel.V(ob.X+p.X, ob.Y+p.Y)))
					prevPoint = p
				}
			}
		}
		if ob.Type == "collider" {
			if len(ob.Ellipses) > 0 {
				colliders = append(colliders, pixel.C(pixel.V(ob.X, ob.Y), ob.Width/2))
			} else {
				colliders = append(colliders, pixel.R(ob.X, ob.Y, ob.X+ob.Width, ob.Y+ob.Height))
			}
		}
		//fmt.Printf("%+v\n", ob)
	}

	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, &startPos, &colliders)

}

func flipY(p *tiled.Object, totalHeight float64) {
	p.Y = totalHeight - p.Y - p.Height
}
func scaleX(p *tiled.Object, scalingFac float64) {
	p.X = p.X * scalingFac
}
func scaleY(p *tiled.Object, scalingFac float64) {
	p.Y = p.Y * scalingFac
}
func PFlipY(o *tiled.Point, totalHeight float64) {
	o.Y = totalHeight - o.Y - 1
}
func PScaleX(o *tiled.Point, scalingFac float64) {
	o.X = scalingFac * o.X
}
func PScaleY(o *tiled.Point, scalingFac float64) {
	o.Y = scalingFac * o.Y
}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
