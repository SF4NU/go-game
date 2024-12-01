package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"go-game/src"
	"go-game/src/assets/sound"
	"go-game/src/assets/sound/effects"
	"go-game/src/assets/sound/ost"
	"go-game/src/resources/camera"
	"go-game/src/resources/monsters/slime"
	"go-game/src/resources/player"
	"image"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

type Game struct {
	Player           *player.Player
	inputSystem      input.System
	Camera           *camera.Camera
	soundtrackPlayer *ost.SoundtrackPlayer
	effectManager    *effects.CombatEffectManager
	Slime            *slime.Slime
}

func (g *Game) Update() error {
	g.Camera.Cam.LookAt(float64(g.Player.Pos.X), float64(g.Player.Pos.Y))
	g.inputSystem.Update()
	g.Player.Update()
	g.Slime.Update()
	g.soundtrackPlayer.Update()
	src.InputHandler()
	if ok := g.Player.Hitbox.IsIntersecting(g.Slime.Hitbox); ok {
		log.Println("INTERSECTING")
	} else {
		log.Println("NOT INTERSECTING")
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var (
		slimo *ebiten.Image
		op    *ebiten.DrawImageOptions
	)
	slimo, op = g.Slime.Get()
	g.Camera.Cam.Draw(slimo, op, screen)
	var playerImg *ebiten.Image
	playerImg, op = g.Player.Get()
	g.Camera.Cam.Draw(playerImg, op, screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func newGame() *Game {
	g := &Game{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})
	ctx := sound.NewSoundContext()
	soundtrackPlayer := ost.NewSoundtrackPlayer(ctx.AudioContext)
	soundtrackPlayer.PlayNext()
	g.soundtrackPlayer = soundtrackPlayer
	effectManager := effects.NewManager(ctx.AudioContext)
	g.effectManager = effectManager
	// todo make a config with all these keys
	keymap := input.Keymap{
		player.ActionMoveLeft:    {input.KeyGamepadLeft, input.KeyA},
		player.ActionMoveRight:   {input.KeyGamepadRight, input.KeyD},
		player.ActionMoveUp:      {input.KeyGamepadUp, input.KeyW},
		player.ActionMoveDown:    {input.KeyGamepadDown, input.KeyS},
		player.ActionItem1:       {input.KeyGamepadR1, input.Key1, input.KeyNum1},
		player.ActionItem2:       {input.KeyGamepadR1, input.Key2, input.KeyNum2},
		player.ActionAttackUp:    {input.KeyGamepadRStickUp, input.KeyUp},
		player.ActionAttackDown:  {input.KeyGamepadRStickDown, input.KeyDown},
		player.ActionAttackLeft:  {input.KeyGamepadRStickLeft, input.KeyLeft},
		player.ActionAttackRight: {input.KeyGamepadRStickRight, input.KeyRight},
		player.ActionJumpUp:      {input.KeyGamepadX, input.KeySpace},
		player.ActionJumpDown:    {input.KeyGamepadX, input.KeySpace},
		player.ActionJumpLeft:    {input.KeyGamepadX, input.KeySpace},
		player.ActionJumpRight:   {input.KeyGamepadX, input.KeySpace},
	}
	g.Player = &player.Player{
		Input: g.inputSystem.NewHandler(0, keymap),
		Pos:   image.Point{X: 512, Y: 384},
	}
	g.Player.PlayerAnimationsNormal().WithEffectMng(effectManager)
	g.Slime = &slime.Slime{
		Pos: image.Point{X: 412, Y: 384},
	}
	g.Slime.PlaySlimeAnimations()
	g.Camera = camera.New(screenWidth, screenHeight, g.Player.Pos.X, g.Player.Pos.Y)
	return g
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebitengine Demo)")
	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
