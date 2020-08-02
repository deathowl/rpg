package world

import (
	"github.com/faiface/pixel"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

func RenderTilemap(tilemap *tiled.Map) *pixel.PictureData {
	renderer, _ := render.NewRenderer(tilemap)
	renderer.RenderVisibleLayers()
	pd := pixel.PictureDataFromImage(renderer.Result)
	renderer.Clear()
	return pd
}
