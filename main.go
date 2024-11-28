package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"go-game/src/assets/player"
	"image"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 480
	screenHeight = 360
)

type Game struct {
	count       int
	Player      *player.Player
	inputSystem input.System
}

func (g *Game) Update() error {
	g.count++
	g.inputSystem.Update()
	g.Player.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)
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

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}
