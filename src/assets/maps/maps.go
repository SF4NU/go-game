package maps

import (
	"fmt"
	"image"
	"os"

	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

const ASSET_PATH = "./assets/map/"

func LoadMap() *image.Rectangle {
	mapString := ASSET_PATH + "example.tmx"
	print(mapString)
	loadedMap, err := tiled.LoadFile(mapString)
	if err != nil {
		panic(err)
	}
	renderer, err := render.NewRenderer(loadedMap)
	if err != nil {
			fmt.Printf("map unsupported for rendering: %s", err.Error())
			os.Exit(2)
	}

	// Render just layer 0 to the Renderer.
	err = renderer.RenderVisibleLayers()
	if err != nil {
			fmt.Printf("layer unsupported for rendering: %s", err.Error())
			os.Exit(2)
	}

	return &renderer.Result.Rect
}