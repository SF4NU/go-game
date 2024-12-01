package sound

import "github.com/hajimehoshi/ebiten/v2/audio"

type Sound struct {
	AudioContext *audio.Context
}

func NewSoundContext() *Sound {
	return &Sound{
		AudioContext: audio.NewContext(44100),
	}
}
