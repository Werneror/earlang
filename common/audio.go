package common

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
)

var initialized = false

func PlayAudio(path string) error {
	audioFile, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "failed to open audio file %s", path)
	}
	defer audioFile.Close()

	audioStreamer, format, err := mp3.Decode(audioFile)
	if err != nil {
		audioStreamer, format, err = wav.Decode(audioFile)
		if err != nil {
			return errors.Wrapf(err, "failed to decode audio file %s", path)
		}
	}
	defer audioStreamer.Close()

	if !initialized {
		initialized = true
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			return errors.Wrap(err, "failed to init speaker")
		}
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(audioStreamer, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
