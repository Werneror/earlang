package common

import (
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var sr = beep.SampleRate(48000)

func init() {
	err := speaker.Init(sr, sr.N(time.Second/2))
	if err != nil {
		logrus.Errorf("failed to init speaker: %v", err)
	}
}

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

	resampled := beep.Resample(4, format.SampleRate, sr, audioStreamer)
	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))
	<-done
	return nil
}
