package player

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/setanarut/anim"
	"github.com/solarlune/resolv"
	"go-game/src/assets/sound/effects"
	"image"
	"io"
	"log"
	"os"
	"time"
)

type Player struct {
	Input              *input.Handler
	Pos                image.Point
	animPlayer         *anim.AnimationPlayer
	class              string
	effectsMng         *effects.CombatEffectManager
	idle               bool
	isPlayingFootsteps bool
	isSwordSwinging    bool
	Hitbox             *resolv.ConvexPolygon
}

var (
	animPlayer      *anim.AnimationPlayer
	classPlayerPath = map[string]string{
		"normal":  "./src/resources/player/sprites/mcNormal.png",
		"sword":   "./src/resources/player/sprites/mcSword.png",
		"staff":   "./src/resources/player/sprites/mcStaff.png",
		"bow":     "./src/resources/player/sprites/mcBow.png",
		"axe":     "./src/resources/player/sprites/mcAxe.png",
		"pickaxe": "./src/resources/player/sprites/mcPickAxe.png",
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
	ActionAttackUp
	ActionAttackDown
	ActionAttackLeft
	ActionAttackRight
	ActionJumpUp
	ActionJumpDown
	ActionJumpLeft
	ActionJumpRight
)

func (p *Player) Update() {
	lastX := p.Pos.X
	lastY := p.Pos.Y
	// basic movement
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
	// jumping
	//if p.Input.ActionIsPressed(ActionJumpUp) && p.Input.ActionIsPressed(ActionMoveUp) {
	//	lockPlayerState(100)
	//	p.animPlayer.SetState("jumpUp")
	//	p.Pos.Y -= 2
	//	time.Sleep(time.Millisecond * 30)
	//	p.Pos.Y -= 2
	//	time.Sleep(time.Millisecond * 30)
	//	p.Pos.Y -= 2
	//	time.Sleep(time.Millisecond * 30)
	//	p.Pos.Y -= 2
	//	time.Sleep(time.Millisecond * 30)
	//	lock = false
	//} else if p.Input.ActionIsJustReleased(ActionJumpUp) {
	//	p.animPlayer.SetState("idleUp")
	//}

	// basic attacks
	if p.Input.ActionIsPressed(ActionAttackUp) && holster[ActionItem1] == true {
		p.Pos.Y -= 1
		p.animPlayer.SetState("attackUp")
	} else if p.Input.ActionIsJustReleased(ActionAttackUp) {
		p.animPlayer.SetState("idleUp")
		p.effectsMng.CloseEffect(1)
	}
	if p.Input.ActionIsPressed(ActionAttackDown) && holster[ActionItem1] == true {
		p.Pos.Y += 1
		p.animPlayer.SetState("attackDown")
	} else if p.Input.ActionIsJustReleased(ActionAttackDown) {
		p.animPlayer.SetState("idleDown")
		p.effectsMng.CloseEffect(1)
	}
	if p.Input.ActionIsPressed(ActionAttackLeft) && holster[ActionItem1] == true {
		p.Pos.X -= 1
		p.animPlayer.SetState("attackLeft")
	} else if p.Input.ActionIsJustReleased(ActionAttackLeft) {
		p.animPlayer.SetState("idleLeft")
		p.effectsMng.CloseEffect(1)
	}
	if p.Input.ActionIsPressed(ActionAttackRight) && holster[ActionItem1] == true {
		p.Pos.X += 1
		p.animPlayer.SetState("attackRight")
	} else if p.Input.ActionIsJustReleased(ActionAttackRight) {
		p.animPlayer.SetState("idleRight")
		p.effectsMng.CloseEffect(1)
	}
	// attacks effect
	if p.Input.ActionIsJustPressed(ActionAttackUp) && holster[ActionItem1] == true {
		p.effectsMng.PlayEffect(1, true)
	}
	if p.Input.ActionIsPressed(ActionAttackDown) && holster[ActionItem1] == true {
		p.effectsMng.PlayEffect(1, true)
	}
	if p.Input.ActionIsPressed(ActionAttackLeft) && holster[ActionItem1] == true {
		p.effectsMng.PlayEffect(1, true)
	}
	if p.Input.ActionIsPressed(ActionAttackRight) && holster[ActionItem1] == true {
		p.effectsMng.PlayEffect(1, true)
	}

	// changing hold items
	if p.Input.ActionIsPressed(ActionItem1) && !lock {
		if holster[ActionItem1] == false {
			p.class = classPlayerPath["sword"]
			p.PlayerAnimationsNormal()
			updateStatusMap(ActionItem1)
			lockPlayerState(200)
		} else {
			p.NormalClass()
			p.PlayerAnimationsNormal()
			holster[ActionItem1] = false
			lockPlayerState(200)
		}
	}
	if p.Input.ActionIsPressed(ActionItem2) && !lock {
		if holster[ActionItem2] == false {
			p.class = classPlayerPath["bow"]
			p.PlayerAnimationsNormal()
			updateStatusMap(ActionItem2)
			lockPlayerState(200)
		} else {
			p.NormalClass()
			p.PlayerAnimationsNormal()
			holster[ActionItem2] = false
			lockPlayerState(200)
		}
	}
	// idle setter
	if p.Pos.X == lastX && p.Pos.Y == lastY {
		p.idle = true
		p.StopFootsteps()
	} else {
		p.idle = false
		p.PlayFootsteps()
	}
	p.Hitbox = resolv.NewRectangle(float64(p.Pos.X), float64(p.Pos.Y), 12, 29)
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
	animPlayer.NewAnimationState("attackUp", 144, 336, 48, 48, 8, true, false)
	animPlayer.NewAnimationState("attackDown", 144, 192, 48, 48, 8, true, false)
	animPlayer.NewAnimationState("attackLeft", 144, 288, 48, 48, 8, true, false)
	animPlayer.NewAnimationState("attackRight", 144, 240, 48, 48, 8, true, false)
	animPlayer.NewAnimationState("jumpDown", 288, 0, 48, 48, 3, false, false)
	animPlayer.NewAnimationState("jumpRight", 288, 48, 48, 48, 3, false, false)
	animPlayer.NewAnimationState("jumpLeft", 288, 96, 48, 48, 3, false, false)
	animPlayer.NewAnimationState("jumpUp", 288, 144, 48, 48, 3, false, false)
	animPlayer.SetState("idleDown")
	animPlayer.SetStateFPS("idleDown", 1)
	animPlayer.SetStateFPS("idleUp", 1)
	animPlayer.SetStateFPS("idleLeft", 1)
	animPlayer.SetStateFPS("idleRight", 1)
	animPlayer.SetStateFPS("walkDown", 6)
	animPlayer.SetStateFPS("walkLeft", 6)
	animPlayer.SetStateFPS("walkRight", 6)
	animPlayer.SetStateFPS("walkUp", 6)
	animPlayer.SetStateFPS("attackUp", 5)
	animPlayer.SetStateFPS("attackDown", 5)
	animPlayer.SetStateFPS("attackLeft", 5)
	animPlayer.SetStateFPS("attackRight", 5)
	animPlayer.SetStateFPS("jumpUp", 3)
	animPlayer.SetStateFPS("jumpDown", 3)
	animPlayer.SetStateFPS("jumpLeft", 3)
	animPlayer.SetStateFPS("jumpRight", 3)
	p.animPlayer = animPlayer
	return p
}

func (p *Player) NormalClass() *Player {
	p.class = classPlayerPath["normal"]
	return p
}

func lockPlayerState(ms time.Duration) {
	go func() {
		lock = true
		time.Sleep(time.Millisecond * ms)
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

func (p *Player) WithEffectMng(effMng *effects.CombatEffectManager) *Player {
	p.effectsMng = effMng
	return p
}

func (p *Player) PlayFootsteps() {
	if !p.idle && !p.isPlayingFootsteps { // Check if idle or already playing
		p.effectsMng.PlayEffect(0, true) // Play the sound effect
		p.isPlayingFootsteps = true      // Set the flag to prevent multiple calls
	}
}

func (p *Player) StopFootsteps() {
	if p.idle && p.isPlayingFootsteps { // Check if idle and currently playing
		p.effectsMng.CloseEffect(0)  // Stop the sound effect
		p.isPlayingFootsteps = false // Reset the flag
	}
}

func (p *Player) PlaySwordSwing() {
	if !p.idle && !p.isPlayingFootsteps { // Check if idle or already playing
		p.effectsMng.PlayEffect(1, true) // Play the sound effect
		p.isPlayingFootsteps = true      // Set the flag to prevent multiple calls
	}
}

func (p *Player) StopSwordSwing() {
	if p.idle && p.isPlayingFootsteps { // Check if idle and currently playing
		p.effectsMng.CloseEffect(1)  // Stop the sound effect
		p.isPlayingFootsteps = false // Reset the flag
	}
}
