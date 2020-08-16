package main

import (
	"fmt"
	"math"
	"time"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/enemy"
	"github.com/deathowl/rpg/engine"
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

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, renderedFg pixel.Picture, initialPos *pixel.Vec, colliders *[]interface{}, enemies *[]*enemy.Enemy) {
	batches := make([]*pixel.Batch, 0)
	var (
		camPos       = *initialPos
		camSpeed     = 40.0
		camZoom      = 4.0
		camZoomSpeed = 1.2
	)

	last := time.Now()
	bgsprite := pixel.NewSprite(renderedBg, renderedBg.Bounds())
	fgsprite := pixel.NewSprite(renderedFg, renderedFg.Bounds())
	mat := pixel.IM
	mat = mat.Moved(win.Bounds().Center())
	mat = mat.ScaledXY(win.Bounds().Center(), pixel.V(win.Bounds().Size().X/renderedBg.Bounds().Size().X, win.Bounds().Size().Y/renderedBg.Bounds().Size().Y))
	sheet, anims, err := engine.LoadAnimationSheet("assets/sheet.png", "assets/spritesheet.csv", 12)
	panicIfErr(err)
	phys := &player.PlayerPhys{
		RunSpeed:  camSpeed,
		Colliders: colliders,
	}
	anim := &player.PlayerAnim{
		Sheet: sheet,
		Anims: anims,
		Rate:  1.0 / 10,
		Dir:   -1,
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
		bgsprite.Draw(win, mat)
		colliderd := imdraw.New(nil)
		colliderd.Push(camPos)
		colliderd.Circle(8, 1)
		for _, collider := range *colliders {
			switch v := collider.(type) {
			case pixel.Circle:
				colliderd.Push(v.Center)
				colliderd.Circle(v.Radius, 1)
			case pixel.Rect:
				colliderd.Push(v.Min, v.Max)
				colliderd.Rectangle(1)
			case pixel.Line:
				colliderd.Push(v.A, v.B)
				colliderd.Line(1)
			}
		}
		anim.Update(dt, phys)
		anim.Draw(win, &camPos)
		fgsprite.Draw(win, mat)
		if win.Pressed(pixelgl.KeyRightControl) {
			colliderd.Draw(win)
		}
		for _, e := range *enemies {
			e.Update(dt, tilemap)
			e.Draw(win)
		}
		frames++
		select {
		case <-second:
			if win.Pressed(pixelgl.KeyRightControl) {

			}
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
	renderedBg := world.RenderBackground(&tilemap)
	renderedFg := world.RenderForeground(&tilemap)

	startPos := win.Bounds().Center()
	scalingFacX := win.Bounds().Size().X / renderedBg.Bounds().Size().X
	scalingFacY := win.Bounds().Size().Y / renderedBg.Bounds().Size().Y
	colliders := make([]interface{}, 0)
	enemies := make([]*enemy.Enemy, 0)
	for _, ob := range tilemap.ObjectGroups[0].Objects {
		if ob.Type != "border" {
			engine.FlipY(ob, renderedBg.Bounds().Size().Y)
			engine.ScaleX(ob, scalingFacX)
			engine.ScaleY(ob, scalingFacY)
		}

		if ob.Type == "start" {
			//startPos = pixel.V(ob.X, math.Abs(ob.Y)).ScaledXY(scalingVec)
			startPos = pixel.Vec{X: ob.X, Y: ob.Y}
		}
		if ob.Type == "border" {
			var prevPoint *tiled.Point
			points := *ob.Polygons[0].Points
			for idx, p := range points {
				if idx == 0 {
					prevPoint = p
				} else if idx == len(points)-1 {
					l1 := pixel.L(pixel.V(ob.X+prevPoint.X, ob.Y+prevPoint.Y), pixel.V(ob.X+p.X, ob.Y+p.Y))
					engine.LFlipY(&l1, renderedBg.Bounds().Size().Y)
					engine.LScaleX(&l1, scalingFacX)
					engine.LScaleY(&l1, scalingFacY)
					l2 := pixel.L(pixel.V(ob.X+p.X, ob.Y+p.Y), pixel.V(ob.X+points[0].X, ob.Y+points[0].Y))
					engine.LFlipY(&l2, renderedBg.Bounds().Size().Y)
					engine.LScaleX(&l2, scalingFacX)
					engine.LScaleY(&l2, scalingFacY)
					colliders = append(colliders, l1, l2)

				} else {
					coll := pixel.L(pixel.V(ob.X+prevPoint.X, ob.Y+prevPoint.Y), pixel.V(ob.X+p.X, ob.Y+p.Y))
					engine.LFlipY(&coll, renderedBg.Bounds().Size().Y)
					engine.LScaleX(&coll, scalingFacX)
					engine.LScaleY(&coll, scalingFacY)
					colliders = append(colliders, coll)
					prevPoint = p
				}
			}
		}
		if ob.Type == "collider" {
			if len(ob.Ellipses) > 0 {
				colliders = append(colliders, pixel.C(pixel.V(ob.X+8, ob.Y+8), ob.Width/2))
			} else {
				colliders = append(colliders, pixel.R(ob.X, ob.Y, ob.X+ob.Width*scalingFacX, ob.Y+ob.Height*scalingFacY))
			}
		}
	}
	for _, eobj := range tilemap.ObjectGroups[1].Objects {
		engine.FlipY(eobj, renderedBg.Bounds().Size().Y)
		engine.ScaleX(eobj, scalingFacX)
		engine.ScaleY(eobj, scalingFacY)
		enemies = append(enemies, enemy.NewEnemy(eobj))
		colliders = append(colliders, pixel.C(pixel.V(eobj.X+8, eobj.Y+8), eobj.Width/2))

	}

	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, renderedFg, &startPos, &colliders, &enemies)

}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
