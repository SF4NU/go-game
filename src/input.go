package src

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

var (
	lock bool
)

// InputHandler handles general input needed for closing game or fullscreen, many others to add...
func InputHandler() {
	fullScreen()
	exitGame()
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

func exitGame() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
}
