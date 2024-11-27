package main

import (
	"bytes"
	input "github.com/quasilyte/ebitengine-input"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
)

const (
	screenWidth  = 480
	screenHeight = 360

	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var (
	runnerImage *ebiten.Image
)

type Game struct {
	count       int
	p           *player
	inputSystem input.System
}

func (g *Game) Update() error {
	g.count++
	g.inputSystem.Update()
	g.p.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	character := runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
	g.p.Draw(screen, character)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newExampleGame() *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
	}
	g.p = &player{
		input: g.inputSystem.NewHandler(0, keymap),
		pos:   image.Point{X: 208, Y: 178},
	}
	return g
}

func main() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}
