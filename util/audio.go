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
}

func NewSound(path string) *SFX {

	sfx, err := loadSound(path)
	if err != nil {
		log.Fatal("Error loading audio: ", err)
	}
	log.Println("Load sound: ", path)

	err = speaker.Init(sfx.format.SampleRate, sfx.format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal("Error initializing audio: ", err)
	}
	//defer speaker.Close()
	log.Printf("Init speaker. SampleRate: %d, BufferSize: %d\n", sfx.format.SampleRate, sfx.format.SampleRate.N(time.Second))

	buffer := beep.NewBuffer(sfx.format)
	buffer.Append(sfx.streamer)
	sfx.streamer.Close()

	log.Println("Len buffer: ", buffer.Len())
	sound1 := buffer.Streamer(0, buffer.Len())
	speaker.Play(sound1)

	// done := make(chan bool)
	// speaker.Play(sound1, beep.Callback(func() {
	// 	done <- true
	// }))
	// <-done

	time.Sleep(1 * time.Second)

	sound1 = buffer.Streamer(0, buffer.Len())
	speaker.Play(sound1)

	// speaker.Play(sound1, beep.Callback(func() {
	// 	done <- true
	// }))
	// <-done

	//	speaker.Play(sound1)

	return sfx
}

func (s *SFX) Close() {
	s.streamer.Close()
}

func (s *SFX) Play() error {
	err := speaker.Init(s.format.SampleRate, s.format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}
	defer speaker.Close()
	done := make(chan bool)
	speaker.Play(beep.Seq(s.streamer, beep.Callback(func() {
		done <- true
	})))
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
