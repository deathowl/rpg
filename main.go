package main

import (
	"fmt"
	"math"
	"time"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/player"
	"github.com/deathowl/rpg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var clearColor = colornames.Black

var (
	frames = 0
	second = time.Tick(time.Second)
)

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, initialPos *pixel.Vec, boundaries *[]pixel.Line) {
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
		RunSpeed: 64,
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
			if camPos.X > camSpeed*dt {
				camPos.X -= camSpeed * dt
			} else {
				fmt.Println(camPos.X)
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			curdir = player.RIGHT

			if camPos.X < win.Bounds().Size().X-camSpeed*dt {
				camPos.X += camSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyDown) {
			curdir = player.DOWN
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			curdir = player.UP
			camPos.Y += camSpeed * dt
		}
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
		// imd.Push(camPos)
		// imd.Circle(3.0, 2.0)
		// imd.Draw(win)
		phys.Update(dt, camPos, curdir)
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
	scalingFacY := renderedBg.Bounds().Size().Y / win.Bounds().Size().Y
	boundaries := make([]pixel.Line, 0)
	for _, ob := range tilemap.ObjectGroups[0].Objects {
		if ob.Type == "start" {
			//fmt.Println(ob.X)
			//fmt.Println(scalingFacX)
			//fmt.Println(ob.Y)
			//fmt.Println(scalingFacY)
			startPos = pixel.Vec{X: ob.X * scalingFacX, Y: ob.Y * scalingFacY}
			//fmt.Println(startPos.X)
			//fmt.Println(startPos.Y)
		}
		if ob.Type == "border" {
			var prevPoint *tiled.Point
			points := *ob.Polygons[0].Points
			for idx, p := range points {
				if idx == 0 {
					prevPoint = p
				} else if idx == len(points)-1 {
					l1 := pixel.L(pixel.V(prevPoint.X, prevPoint.Y), pixel.V(p.X, p.Y))
					l2 := pixel.L(pixel.V(p.X, p.Y), pixel.V(points[0].X, points[0].Y))
					boundaries = append(boundaries, l1, l2)

				} else {
					boundaries = append(boundaries, pixel.L(pixel.V(prevPoint.X, prevPoint.Y), pixel.V(p.X, p.Y)))
					prevPoint = p
				}
				fmt.Println(boundaries)
			}
		}
		if ob.Type == "collider" {
			fmt.Printf("%+v\n", ob)
		}
		//fmt.Printf("%+v\n", ob)
	}

	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, &startPos, &boundaries)

}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
