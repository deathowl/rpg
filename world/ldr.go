package world
import 
	("github.com/lafriks/go-tiled")


func LoadTileMap(path string) tiled.Map {
	tilemap, err := tiled.LoadFromFile("./island.tmx")
	if (err != nil) {
		panic(err)
	}
	return *tilemap
}