package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

var (
	lock bool
)

func InputHandler() {
	fullScreen()
}

func fullScreen() {
	if ebiten.IsKeyPressed(ebiten.KeyF11) {
		if ebiten.IsFullscreen() && !lock {
			ebiten.SetFullscreen(false)
			lock = true
			go func() {
				time.Sleep(time.Second * 1)
				lock = false
			}()
		} else if !lock {
			ebiten.SetFullscreen(true)
			lock = true
			go func() {
				time.Sleep(time.Second * 1)
				lock = false
			}()
		}
	}
}
