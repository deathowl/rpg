package world

import (
	"github.com/deathowl/go-tiled"
	"github.com/deathowl/go-tiled/render"
	"github.com/faiface/pixel"
)

func RenderTilemap(tilemap *tiled.Map) *pixel.PictureData {
	renderer, _ := render.NewRenderer(tilemap)
	renderer.RenderVisibleLayers()
	pd := pixel.PictureDataFromImage(renderer.Result)
	renderer.Clear()
	return pd
}
