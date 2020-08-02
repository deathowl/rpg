package main

import (
	"fmt"
	"image/png"
	"math"
	"os"
	"time"

	"golang.org/x/image/colornames"

	"github.com/deathowl/rpg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lafriks/go-tiled"
)

var clearColor = colornames.Skyblue

var (
	frames = 0
	second = time.Tick(time.Second)
)

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, initialPos *pixel.Vec) {
	batches := make([]*pixel.Batch, 0)

	var (
		camPos       = *initialPos
		camSpeed     = 1000.0
		camZoom      = 4.0
		camZoomSpeed = 1.2
	)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Camera movement
		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		win.Clear(clearColor)

		// Draw tiles
		for _, batch := range batches {
			batch.Clear()
		}

		rsprirte := pixel.NewSprite(renderedBg, renderedBg.Bounds())
		mat := pixel.IM
		rsprirte.Draw(win, mat)
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
	startPos := &pixel.ZV
	for _, ob := range tilemap.ObjectGroups[0].Objects {
		fmt.Println(ob.Name)
		if ob.Type == "start" {
			startPos = &pixel.Vec{X: ob.X / 2, Y: ob.Y}
		}
		fmt.Printf("%+v\n", ob)
	}

	panicIfErr(err)

	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, startPos)
}

func loadSprite(path string) (*pixel.Sprite, *pixel.PictureData) {
	fmt.Println(path)
	f, err := os.Open(path)
	panicIfErr(err)

	img, err := png.Decode(f)
	panicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds()), pd
}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
