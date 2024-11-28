package maps

import (
	"image"
	"os"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

type MapRenderer struct {
	mapData *tiled.Map
	tileCache map[uint32]*ebiten.Image
}

func NewMapRenderer(mapFilePath string) (*MapRenderer, error) {
	mapData, err := tiled.LoadFile(mapFilePath)
	if err != nil {
		return nil, err
	}

	return &MapRenderer{
		mapData: mapData,
		tileCache: make(map[uint32]*ebiten.Image),
	}, nil
}

func (mr *MapRenderer) Draw(screen *ebiten.Image) {
	for _, layer := range mr.mapData.Layers {
		var xs, xe, xi, ys, ye, yi int
		xs = 0
		xe = mr.mapData.Width
		xi = 1
		ys = 0
		ye = mr.mapData.Height
		yi = 1

		i := 0
		for y := ys; y*yi < ye; y = y + yi {
			for x := xs; x*xi < xe; x = x + xi {
				if layer.Tiles[i].IsNil() {
					i++
					continue
				}

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x * mr.mapData.TileWidth), float64(y * mr.mapData.TileHeight))

				cached_tile, ok := mr.tileCache[layer.Tiles[i].Tileset.FirstGID + layer.Tiles[i].ID]
				if ok {
					screen.DrawImage(cached_tile, op)
				}

				i++
			}
		}
	}
}

func (e *MapRenderer) GetTilePosition(x, y int) image.Rectangle {
	return image.Rect(x*e.mapData.TileWidth,
		y*e.mapData.TileHeight,
		(x+1)*e.mapData.TileWidth,
		(y+1)*e.mapData.TileHeight)
}


func (e *MapRenderer) RotateTileImage(tile *tiled.LayerTile, img image.Image) image.Image {
	timg := img
	if tile.HorizontalFlip {
		timg = imaging.FlipH(timg)
	}
	if tile.VerticalFlip {
		timg = imaging.FlipV(timg)
	}
	if tile.DiagonalFlip {
		timg = imaging.FlipH(imaging.Rotate90(timg))
	}

	return timg
}

func (mr *MapRenderer) PreloadTiles() error {
	for _, tileset := range mr.mapData.Tilesets {
		file, err := os.Open(tileset.GetFileFullPath(tileset.Image.Source))
		if err != nil {
			return err
		}
		defer file.Close()

		// Load the full tileset image
		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}
		tilesetImage := ebiten.NewImageFromImage(img)

		// Crop and store each tile as an ebiten.Image
		for i := uint32(0); i < uint32(tileset.TileCount); i++ {
			rect := tileset.GetTileRect(i)
			tileImg := tilesetImage.SubImage(rect).(*ebiten.Image)
			mr.tileCache[i+tileset.FirstGID] = tileImg
		}
	}
	return nil
}