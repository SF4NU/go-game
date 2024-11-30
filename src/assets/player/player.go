package player

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/setanarut/anim"
	"image"
	"io"
	"log"
	"os"
)

var animPlayer *anim.AnimationPlayer

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
)

type Player struct {
	Input      *input.Handler
	Pos        image.Point
	animPlayer *anim.AnimationPlayer
}

func (p *Player) Update() {
	if p.Input.ActionIsPressed(ActionMoveLeft) {
		p.Pos.X -= 2
		p.animPlayer.SetState("moveLeft")
	} else if p.Input.ActionIsJustReleased(ActionMoveLeft) {
		p.animPlayer.SetState("idleLeft")
	}
	if p.Input.ActionIsPressed(ActionMoveRight) {
		p.Pos.X += 2
		p.animPlayer.SetState("moveRight")
	} else if p.Input.ActionIsJustReleased(ActionMoveRight) {
		p.animPlayer.SetState("idleRight")
	}
	if p.Input.ActionIsPressed(ActionMoveUp) {
		p.Pos.Y -= 2
		p.animPlayer.SetState("moveUp")
	} else if p.Input.ActionIsJustReleased(ActionMoveUp) {
		p.animPlayer.SetState("idleUp")
	}
	if p.Input.ActionIsPressed(ActionMoveDown) {
		p.Pos.Y += 2
		p.animPlayer.SetState("moveDown")
	} else if p.Input.ActionIsJustReleased(ActionMoveDown) {
		p.animPlayer.SetState("idleDown")
	}

	log.Printf("PosX:%d PosY:%d", p.Pos.X, p.Pos.Y)
	p.animPlayer.Update()
}

func (p *Player) Get() (*ebiten.Image, *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.Pos.X), float64(p.Pos.Y))
	//screen.DrawImage(p.animPlayer.CurrentFrame, op)
	return p.animPlayer.CurrentFrame, op
}

func (p *Player) PlayerAnimations() *Player {
	filePath := "./src/assets/player/character.png"
	// Decode an image from the image file's byte slice.
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	newPlayerImg := ebiten.NewImageFromImage(img)

	animPlayer = anim.NewAnimationPlayer(newPlayerImg)
	animPlayer.NewAnimationState("idleDown", 48, 0, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleUp", 48, 144, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleLeft", 48, 96, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleRight", 48, 48, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("moveDown", 0, 0, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("moveRight", 0, 48, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("moveLeft", 0, 96, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("moveUp", 0, 144, 48, 48, 4, true, false)
	animPlayer.SetState("idleDown")
	animPlayer.SetStateFPS("idleDown", 1)
	animPlayer.SetStateFPS("idleUp", 1)
	animPlayer.SetStateFPS("idleLeft", 1)
	animPlayer.SetStateFPS("idleRight", 1)
	animPlayer.SetStateFPS("moveDown", 6)
	animPlayer.SetStateFPS("moveLeft", 6)
	animPlayer.SetStateFPS("moveRight", 6)
	animPlayer.SetStateFPS("moveUp", 6)
	p.animPlayer = animPlayer
	return p
}
