package ost

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"log"
)

//go:embed tracks/*.ogg
var audioFiles embed.FS

type SoundtrackPlayer struct {
	audioContext *audio.Context
	currentIndex int
	tracks       [][]byte
	player       *audio.Player
}

func NewSoundtrackPlayer(audioContext *audio.Context) *SoundtrackPlayer {
	trackFiles := []string{"tracks/The-Quest-Unfolds.ogg"}
	var tracks [][]byte
	for _, file := range trackFiles {
		data, err := audioFiles.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read sound file: %v", err)
		}
		tracks = append(tracks, data)
	}
	return &SoundtrackPlayer{
		audioContext: audioContext,
		currentIndex: 0,
		tracks:       tracks,
	}
}

func (s *SoundtrackPlayer) PlayNext() {
	if s.player != nil {
		s.player.Close() // Ensure the previous player is cleaned up
	}

	if s.currentIndex >= len(s.tracks) {
		log.Println("No more tracks to play.")
		return
	}

	track := s.tracks[s.currentIndex]
	stream, err := vorbis.DecodeWithSampleRate(44100, bytes.NewReader(track))
	if err != nil {
		log.Fatalf("Failed to decode sound track: %v", err)
	}
	s.player, err = (*s.audioContext).NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create sound player: %v", err)
	}

	s.player.Play()
	s.player.SetVolume(0.05)
	s.currentIndex++
}

func (s *SoundtrackPlayer) Update() {
	if s.player == nil {
		return
	}

	// Check if the current track has finished
	if !s.player.IsPlaying() {
		s.PlayNext()
	}
}
