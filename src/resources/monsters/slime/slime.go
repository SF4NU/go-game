package slime

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/setanarut/anim"
	"github.com/solarlune/resolv"
	"go-game/src/assets/sound/effects"
	"image"
	"io"
	"log"
	"os"
)

type Slime struct {
	Pos        image.Point
	animSlime  *anim.AnimationPlayer
	effectsMng *effects.CombatEffectManager
	idle       bool
	isHurt     bool
	Hitbox     *resolv.Circle
}

func (s *Slime) Update() {
	s.Hitbox = resolv.NewCircle(float64(s.Pos.X), float64(s.Pos.Y), 14)
	s.animSlime.Update()
}

func (s *Slime) Get() (*ebiten.Image, *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.Pos.X), float64(s.Pos.Y))
	return s.animSlime.CurrentFrame, op
}

func (s *Slime) ChangeSprite() *ebiten.Image {
	file, err := os.Open("./src/resources/monsters/slime/sprites/slime.png")
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
	newSlimeImg := ebiten.NewImageFromImage(img)
	return newSlimeImg
}

func (s *Slime) PlaySlimeAnimations() *Slime {
	animSlime := anim.NewAnimationPlayer(s.ChangeSprite())
	animSlime.NewAnimationState("idle", 0, 0, 32, 25, 8, true, false)
	animSlime.SetState("idle")
	animSlime.SetStateFPS("idle", 6)
	s.animSlime = animSlime
	return s
}
