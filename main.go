package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
	input "github.com/quasilyte/ebitengine-input"
	"go-game/src/assets/player"
	"image"
	_ "image/png"
	"log"
	"os"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
)

type Game struct {
	count       int
	Player      *player.Player
	inputSystem input.System
	mapImg      *image.NRGBA
}

func (g *Game) Update() error {
	g.count++
	g.inputSystem.Update()
	g.Player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenImg := ebiten.NewImageFromImage(g.mapImg)
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(ebitenImg, op)
	g.Player.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newExampleGame() *Game {
	gameMap, err := tiled.LoadFile(mapPath)
	if err != nil {
		fmt.Printf("error parsing map: %s", err.Error())
		os.Exit(2)
	}
	renderer, err := render.NewRenderer(gameMap)
	if err != nil {
		fmt.Printf("map unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}
	err = renderer.RenderVisibleLayers()
	if err != nil {
		fmt.Printf("layer unsupported for rendering: %s", err.Error())
		os.Exit(2)
	}
	img := renderer.Result
	renderer.Clear()
	g := &Game{
		mapImg: img,
	}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	keymap := input.Keymap{
		player.ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		player.ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		player.ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		player.ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
	}
	g.Player = &player.Player{
		Input: g.inputSystem.NewHandler(0, keymap),
		Pos:   image.Point{X: 208, Y: 178},
	}
	g.Player.PlayerAnimations()
	return g
}

const mapPath = "./src/assets/maps/assets/MagicLand.tmx"

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}
