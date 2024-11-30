package audio

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"io"
	"log"
)

type audioStream interface {
	io.ReadSeeker
	Length() int64
}

const (
	sampleRate     = 44100
	bytesPerSample = 8
)

//go:embed ost/*.ogg
var audioFiles embed.FS

type SoundtrackPlayer struct {
	audioContext *audio.Context
	currentIndex int
	tracks       [][]byte
	player       *audio.Player
}

func NewSoundtrackPlayer() *SoundtrackPlayer {
	audioContext := audio.NewContext(sampleRate)
	trackFiles := []string{"ost/The-Quest-Unfolds.ogg"}
	var tracks [][]byte
	for _, file := range trackFiles {
		data, err := audioFiles.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read audio file: %v", err)
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
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(track))
	if err != nil {
		log.Fatalf("Failed to decode audio track: %v", err)
	}
	s.player, err = (*s.audioContext).NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create audio player: %v", err)
	}

	s.player.Play()
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
