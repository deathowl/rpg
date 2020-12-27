package world

import (
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	"github.com/faiface/pixel"
)

func RenderBackground(tilemap *tiled.Map) *pixel.PictureData {
	renderer, _ := render.NewRenderer(tilemap)
	renderer.RenderLayer(0)
	renderer.RenderLayer(1)
	pd := pixel.PictureDataFromImage(renderer.Result)
	renderer.Clear()
	return pd
}

func RenderForeground(tilemap *tiled.Map) *pixel.PictureData {
	renderer, _ := render.NewRenderer(tilemap)
	renderer.RenderLayer(2)
	renderer.RenderLayer(3)
	pd := pixel.PictureDataFromImage(renderer.Result)
	renderer.Clear()
	return pd
}
