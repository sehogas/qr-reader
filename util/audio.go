package util

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type SFX struct {
	streamer beep.StreamSeekCloser
	format   beep.Format
	buffer   *beep.Buffer
}

func NewSound(path string) *SFX {

	sfx, err := loadSound(path)
	if err != nil {
		log.Fatal("Error loading audio: ", err)
	}

	sfx.buffer = beep.NewBuffer(sfx.format)
	sfx.buffer.Append(sfx.streamer)
	sfx.streamer.Close()

	return sfx
}

func (s *SFX) Play() error {
	err := speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}

	sound := s.buffer.Streamer(0, s.buffer.Len())

	done := make(chan bool)
	speaker.Play(sound, beep.Callback(func() {
		done <- true
	}))
	<-done

	return nil
}

func loadSound(path string) (*SFX, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}

	return &SFX{
		streamer: streamer,
		format:   format,
	}, nil
}
