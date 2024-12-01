package effects

import (
	"bytes"
	"embed"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

//go:embed tracks/*.mp3
var audioFiles embed.FS

type effect struct {
	Name       string
	Sound      *audio.Player // Sound effect associated with the effect
	IsLooping  bool          // Indicates whether this effect should loop
	playingMux sync.Mutex    // Mutex to ensure thread-safe state changes
}

type CombatEffectManager struct {
	audioContext *audio.Context
	effects      map[int]*effect
	closed       bool
}

func NewManager(audioContext *audio.Context) *CombatEffectManager {
	newCombatEffectManager := &CombatEffectManager{
		audioContext: audioContext,
		effects:      make(map[int]*effect, 12),
	}
	err := newCombatEffectManager.LoadSoundEffect()
	if err != nil {
		panic(err)
	}
	return newCombatEffectManager
}

var trackFiles = []string{"tracks/grass-normal.mp3", "tracks/sword-slash-1.mp3"}

func (m *CombatEffectManager) LoadSoundEffect() error {
	for i, f := range trackFiles {
		track, err := audioFiles.ReadFile(f)
		if err != nil {
			return err
		}

		stream, err := mp3.DecodeWithSampleRate(m.audioContext.SampleRate(), bytes.NewReader(track))
		if err != nil {
			return err
		}

		player, err := m.audioContext.NewPlayer(stream)
		if err != nil {
			return err
		}

		m.effects[i] = &effect{
			Name:      trackFiles[i],
			Sound:     player,
			IsLooping: false, // Default is not looping
		}
	}
	return nil
}

func (m *CombatEffectManager) PlayEffect(num int, loop bool) {
	eff, ok := m.effects[num]
	if !ok {
		panic(num)
	}

	eff.playingMux.Lock()
	defer eff.playingMux.Unlock()

	if eff.Sound.IsPlaying() {
		return // Already playing
	}

	eff.IsLooping = loop

	// Play sound and manage looping
	go func(eff *effect) {
		for {
			eff.Sound.Rewind()
			eff.Sound.SetVolume(0.10)
			eff.Sound.Play()

			// Wait for the sound to finish
			for eff.Sound.IsPlaying() {
				continue
			}

			// Break if no longer looping
			if !eff.IsLooping {
				break
			}
		}
	}(eff)
}

func (m *CombatEffectManager) StopEffect(num int) {
	eff, ok := m.effects[num]
	if !ok {
		panic(num)
	}

	eff.playingMux.Lock()
	defer eff.playingMux.Unlock()

	eff.IsLooping = false // Ensure looping stops
	if eff.Sound.IsPlaying() {
		eff.Sound.Close() // Stop the sound
	}
}

func (m *CombatEffectManager) CloseEffect(num int) {
	eff, ok := m.effects[num]
	if !ok {
		panic(num)
	}

	eff.playingMux.Lock()
	defer eff.playingMux.Unlock()

	eff.IsLooping = false // Stop looping if active
	err := eff.Sound.Close()
	if err != nil {
		panic(err)
	}
}
