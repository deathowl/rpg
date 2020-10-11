package main

import (
	"fmt"
	"math"
	"time"

	"github.com/deathowl/go-tiled"
	"github.com/deathowl/rpg/enemy"
	"github.com/deathowl/rpg/engine"
	"github.com/deathowl/rpg/npc"
	"github.com/deathowl/rpg/player"
	"github.com/deathowl/rpg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var clearColor = colornames.Black

var (
	second = time.Tick(time.Second)
)

func gameloop(win *pixelgl.Window, tilemap *tiled.Map, renderedBg pixel.Picture, renderedFg pixel.Picture, initialPos *pixel.Vec, colliders *[]interface{}, enemies *[]*enemy.Enemy, npcs *[]*npc.NPC) {
	batches := make([]*pixel.Batch, 0)
	var (
		playerPos    = *initialPos
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
	fpsc := engine.NewFPSWatchSimple()
	fpsc.Start()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Camera movement
		cam := pixel.IM.Scaled(playerPos, camZoom).Moved(win.Bounds().Center().Sub(playerPos))
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
		playerPos = phys.Update(dt, playerPos, curdir)
		if win.Pressed(pixelgl.KeySpace) {
			fmt.Println(playerPos.X)
			fmt.Println(playerPos.Y)
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		// if camZoom < 3.5 {
		// 	camZoom = 3.5
		// }
		if camZoom > 8.25 {
			camZoom = 8.25
		}
		win.Clear(clearColor)

		// Draw tiles
		for _, batch := range batches {
			batch.Clear()
		}
		bgsprite.Draw(win, mat)
		colliderd := imdraw.New(nil)
		colliderd.Push(playerPos)
		colliderd.Circle(8, 1)
		for _, collider := range *colliders {
			switch v := collider.(type) {
			case *pixel.Circle:
				colliderd.Push(v.Center)
				colliderd.Circle(v.Radius, 1)
			case *pixel.Rect:
				colliderd.Push(v.Min, v.Max)
				colliderd.Rectangle(1)
			case *pixel.Line:
				colliderd.Push(v.A, v.B)
				colliderd.Line(1)
			}
		}
		anim.Update(dt, phys)
		anim.Draw(win, &playerPos)
		fgsprite.Draw(win, mat)

		// txt := engine.DrawText(pixel.V(playerPos.X-190, playerPos.Y-118), "HEALTH", colornames.White)
		// txt.Draw(win, pixel.IM)
		ebardrawer := imdraw.New(nil)
		for _, e := range *enemies {
			e.Update(dt, colliders, &playerPos)
			ebardrawer.Color = colornames.Red
			ebardrawer.Push(pixel.V((e.Pos.X-5), (e.Pos.Y+8)), pixel.V(e.Pos.X+5, e.Pos.Y+10))
			ebardrawer.Rectangle(0)
			ebardrawer.Draw(win)
			e.Draw(win)
		}
		ebardrawer.Draw(win)

		for _, npc := range *npcs {
			npc.Update(dt)
			// txt := engine.DrawText(npc.Pos, npc.Name, colornames.Black)
			// txt.DrawColorMask(win, pixel.IM, colornames.Black)
			npc.Draw(win)
		}
		if win.Pressed(pixelgl.KeyRightControl) {
			colliderd.Draw(win)
		}
		bardrawer := imdraw.New(nil)
		bardrawer.Push(pixel.V((playerPos.X-238), (playerPos.Y-120)), pixel.V((playerPos.X-178), playerPos.Y-115))
		bardrawer.Rectangle(1)
		bardrawer.Draw(win)
		bardrawer.Color = colornames.Green
		bardrawer.Push(pixel.V((playerPos.X-237), (playerPos.Y-119)), pixel.V(playerPos.X-179, playerPos.Y-116))
		bardrawer.Rectangle(0)
		bardrawer.Draw(win)
		fpsc.Poll()
		fpsc.Draw(win, pixel.V(playerPos.X+100, playerPos.Y+100))
		win.Update()
	}
}

func initialize() {
	fmt.Println("initialize called")
	// Create the window with OpenGL
	cfg := pixelgl.WindowConfig{
		Title:  "Tiled Rpg",
		Bounds: pixel.R(0, 0, 1920, 1080),
		VSync:  false,
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
	npcs := make([]*npc.NPC, 0)
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
					colliders = append(colliders, &l1, &l2)

				} else {
					coll := pixel.L(pixel.V(ob.X+prevPoint.X, ob.Y+prevPoint.Y), pixel.V(ob.X+p.X, ob.Y+p.Y))
					engine.LFlipY(&coll, renderedBg.Bounds().Size().Y)
					engine.LScaleX(&coll, scalingFacX)
					engine.LScaleY(&coll, scalingFacY)
					colliders = append(colliders, &coll)
					prevPoint = p
				}
			}
		}
		if ob.Type == "collider" {
			if len(ob.Ellipses) > 0 {
				coll := pixel.C(pixel.V(ob.X+8, ob.Y+8), ob.Width/2)
				colliders = append(colliders, &coll)
			} else {
				coll := pixel.R(ob.X, ob.Y, ob.X+ob.Width*scalingFacX, ob.Y+ob.Height*scalingFacY)
				colliders = append(colliders, &coll)
			}
		}
	}
	for _, eobj := range tilemap.ObjectGroups[1].Objects {
		engine.FlipY(eobj, renderedBg.Bounds().Size().Y)
		engine.ScaleX(eobj, scalingFacX)
		engine.ScaleY(eobj, scalingFacY)
		collider := pixel.C(pixel.V(eobj.X+8, eobj.Y+8), eobj.Width/2)
		e := enemy.NewEnemy(eobj, &collider)
		colliders = append(colliders, &collider)
		enemies = append(enemies, e)

	}
	for _, eobj := range tilemap.ObjectGroups[2].Objects {
		engine.FlipY(eobj, renderedBg.Bounds().Size().Y)
		engine.ScaleX(eobj, scalingFacX)
		engine.ScaleY(eobj, scalingFacY)
		collider := pixel.C(pixel.V(eobj.X+8, eobj.Y+8), eobj.Width/2)
		e := npc.NewNPC(eobj, &collider)
		colliders = append(colliders, &collider)
		npcs = append(npcs, e)

	}
	fmt.Println("use WASD to move camera around")
	gameloop(win, &tilemap, renderedBg, renderedFg, &startPos, &colliders, &enemies, &npcs)

}

func main() {
	pixelgl.Run(initialize)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
