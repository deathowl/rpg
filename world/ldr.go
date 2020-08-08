package world

import (
	"github.com/deathowl/go-tiled"
)

func LoadTileMap(path string) tiled.Map {
	tilemap, err := tiled.LoadFromFile(path)
	if err != nil {
		panic(err)
	}
	return *tilemap
}
