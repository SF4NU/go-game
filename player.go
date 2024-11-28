package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
)

type player struct {
	input     *input.Handler
	pos       image.Point
	character *ebiten.Image
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.pos.X -= 4
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.pos.X += 4
	}
	if p.input.ActionIsPressed(ActionMoveUp) {
		p.pos.Y -= 4
	}
	if p.input.ActionIsPressed(ActionMoveDown) {
		p.pos.Y += 4
	}
}

//func (p *player) Draw(screen *ebiten.Image) {
//	ebitenutil.DebugPrintAt(screen, "player", p.pos.X, p.pos.Y)
//}

func (p *player) Draw(screen *ebiten.Image, character *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.pos.X), float64(p.pos.Y))
	screen.DrawImage(character, op)
}

