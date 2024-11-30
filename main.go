package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"go-game/src"
	"go-game/src/assets/audio"
	"go-game/src/assets/player"
	"go-game/src/resources/camera"
	"image"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

var (
	lock bool
)

type Game struct {
	Player           *player.Player
	inputSystem      input.System
	Camera           *camera.Camera
	soundtrackPlayer *audio.SoundtrackPlayer
}

func (g *Game) Update() error {
	g.Camera.Cam.LookAt(float64(g.Player.Pos.X), float64(g.Player.Pos.Y))
	g.inputSystem.Update()
	g.Player.Update()
	g.soundtrackPlayer.Update()
	src.InputHandler()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	playerImg, op := g.Player.Get()
	g.Camera.Cam.Draw(playerImg, op, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newExampleGame() *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	newSoundPlayer := audio.NewSoundtrackPlayer()
	newSoundPlayer.PlayNext()
	g.soundtrackPlayer = newSoundPlayer
	keymap := input.Keymap{
		player.ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		player.ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
		player.ActionMoveUp:    {input.KeyGamepadUp, input.KeyUp, input.KeyW},
		player.ActionMoveDown:  {input.KeyGamepadDown, input.KeyDown, input.KeyS},
	}
	g.Player = &player.Player{
		Input: g.inputSystem.NewHandler(0, keymap),
		Pos:   image.Point{X: screenWidth / 2, Y: screenHeight / 2},
	}
	g.Player.PlayerAnimations()
	g.Camera = camera.New(screenWidth, screenHeight, g.Player.Pos.X, g.Player.Pos.Y)
	return g
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}
