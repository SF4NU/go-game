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
	"time"
)

var (
	animPlayer      *anim.AnimationPlayer
	classPlayerPath = map[string]string{
		"normal":  "./src/assets/player/sprites/mcNormal.png",
		"sword":   "./src/assets/player/sprites/mcSword.png",
		"staff":   "./src/assets/player/sprites/mcStaff.png",
		"bow":     "./src/assets/player/sprites/mcBow.png",
		"axe":     "./src/assets/player/sprites/mcAxe.png",
		"pickaxe": "./src/assets/player/sprites/mcPickAxe.png",
	}
	holster = map[input.Action]bool{
		ActionItem1: false,
		ActionItem2: false,
		ActionItem3: false,
		ActionItem4: false,
		ActionItem5: false,
		ActionItem6: false,
		ActionItem7: false,
		ActionItem8: false,
		ActionItem9: false,
	}
	lock bool
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionItem1
	ActionItem2
	ActionItem3
	ActionItem4
	ActionItem5
	ActionItem6
	ActionItem7
	ActionItem8
	ActionItem9
)

type Player struct {
	Input      *input.Handler
	Pos        image.Point
	animPlayer *anim.AnimationPlayer
	class      string
}

func (p *Player) Update() {
	if p.Input.ActionIsPressed(ActionMoveLeft) {
		p.Pos.X -= 2
		p.animPlayer.SetState("walkLeft")
	} else if p.Input.ActionIsJustReleased(ActionMoveLeft) {
		p.animPlayer.SetState("idleLeft")
	}
	if p.Input.ActionIsPressed(ActionMoveRight) {
		p.Pos.X += 2
		p.animPlayer.SetState("walkRight")
	} else if p.Input.ActionIsJustReleased(ActionMoveRight) {
		p.animPlayer.SetState("idleRight")
	}
	if p.Input.ActionIsPressed(ActionMoveUp) {
		p.Pos.Y -= 2
		p.animPlayer.SetState("walkUp")
	} else if p.Input.ActionIsJustReleased(ActionMoveUp) {
		p.animPlayer.SetState("idleUp")
	}
	if p.Input.ActionIsPressed(ActionMoveDown) {
		p.Pos.Y += 2
		p.animPlayer.SetState("walkDown")
	} else if p.Input.ActionIsJustReleased(ActionMoveDown) {
		p.animPlayer.SetState("idleDown")
	}
	if p.Input.ActionIsPressed(ActionItem1) && !lock {
		if holster[ActionItem1] == false {
			p.class = classPlayerPath["sword"]
			p.PlayerAnimationsNormal()
			updateStatusMap(ActionItem1)
			lockClassChange()
		} else {
			p.NormalClass()
			p.PlayerAnimationsNormal()
			holster[ActionItem1] = false
			lockClassChange()
		}
	}
	if p.Input.ActionIsPressed(ActionItem2) && !lock {
		if holster[ActionItem2] == false {
			p.class = classPlayerPath["bow"]
			p.PlayerAnimationsNormal()
			updateStatusMap(ActionItem2)
			lockClassChange()
		} else {
			p.NormalClass()
			p.PlayerAnimationsNormal()
			holster[ActionItem2] = false
			lockClassChange()
		}
	}

	//log.Printf("PosX:%d PosY:%d", p.Pos.X, p.Pos.Y)
	p.animPlayer.Update()
}

func (p *Player) ChangeSprite() *ebiten.Image {
	// Decode an image from the image file's byte slice.
	//p.class = classPlayerPath["normal"]
	if p.class == "" {
		p.class = classPlayerPath["normal"]
	}
	file, err := os.Open(p.class)
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
	return newPlayerImg
}

func (p *Player) Get() (*ebiten.Image, *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.Pos.X), float64(p.Pos.Y))
	//screen.DrawImage(p.animPlayer.CurrentFrame, op)
	return p.animPlayer.CurrentFrame, op
}

func (p *Player) PlayerAnimationsNormal() *Player {
	animPlayer = anim.NewAnimationPlayer(p.ChangeSprite())
	animPlayer.NewAnimationState("idleDown", 48, 0, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleUp", 48, 144, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleLeft", 48, 96, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("idleRight", 48, 48, 48, 48, 1, true, false)
	animPlayer.NewAnimationState("walkDown", 0, 0, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("walkRight", 0, 48, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("walkLeft", 0, 96, 48, 48, 4, true, false)
	animPlayer.NewAnimationState("walkUp", 0, 144, 48, 48, 4, true, false)
	animPlayer.SetState("idleDown")
	animPlayer.SetStateFPS("idleDown", 1)
	animPlayer.SetStateFPS("idleUp", 1)
	animPlayer.SetStateFPS("idleLeft", 1)
	animPlayer.SetStateFPS("idleRight", 1)
	animPlayer.SetStateFPS("walkDown", 6)
	animPlayer.SetStateFPS("walkLeft", 6)
	animPlayer.SetStateFPS("walkRight", 6)
	animPlayer.SetStateFPS("walkUp", 6)
	p.animPlayer = animPlayer
	return p
}

func (p *Player) NormalClass() *Player {
	p.class = classPlayerPath["normal"]
	return p
}

func lockClassChange() {
	go func() {
		lock = true
		time.Sleep(time.Millisecond * 200)
		lock = false
	}()
}

func updateStatusMap(action input.Action) {
	for key := range holster {
		holster[key] = false
	}
	if _, exists := holster[action]; exists {
		holster[action] = true
	}
}
