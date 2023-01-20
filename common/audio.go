package common

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var initialized = false

func PlayMP3(path string) error {
	audioFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer audioFile.Close()

	audioStreamer, format, err := mp3.Decode(audioFile)
	if err != nil {
		return err
	}
	defer audioStreamer.Close()

	if !initialized {
		initialized = true
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return nil
		}
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
