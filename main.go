package main

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/image/colornames"

	"github.com/deathowl/rpg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lafriks/go-tiled"
)

var clearColor = colornames.Black

var (
	frames = 0
	second = time.Tick(time.Second)
)

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, initialPos *pixel.Vec) {
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
	fmt.Println(mat)
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		if win.Pressed(pixelgl.KeyLeft) {
			if camPos.X > camSpeed*dt {
				camPos.X -= camSpeed * dt
			} else {
				fmt.Println(camPos.X)
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if camPos.X < win.Bounds().Size().X-camSpeed*dt {
				camPos.X += camSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
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
		imd := imdraw.New(nil)
		imd.Push(camPos)
		imd.Circle(3.0, 2.0)
		imd.Draw(win)
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
	tilemap := world.LoadTileMap("./island.tmx")
	fmt.Println(tilemap)
	renderedBg := world.RenderTilemap(&tilemap)

	startPos := win.Bounds().Center()
	scalingFacX := win.Bounds().Size().X / renderedBg.Bounds().Size().X
	scalingFacY := renderedBg.Bounds().Size().Y / win.Bounds().Size().Y
	for _, ob := range tilemap.ObjectGroups[0].Objects {
		fmt.Println(ob.Name)
		if ob.Type == "start" {
			fmt.Println(ob.X)
			fmt.Println(scalingFacX)
			fmt.Println(ob.Y)
			fmt.Println(scalingFacY)
			startPos = pixel.Vec{X: ob.X * scalingFacX, Y: ob.Y * scalingFacY}
			fmt.Println(startPos.X)
			fmt.Println(startPos.Y)
		}
		fmt.Printf("%+v\n", ob)
	}

	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, &startPos)

}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
