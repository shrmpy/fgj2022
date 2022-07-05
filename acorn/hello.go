//go:build !(wasm || js || android)

package acorn

import (
	_ "embed"

	//"log"
)
import (
	"github.com/gen2brain/flite-go"
)

type Parcel struct {
	voice   *flite.Voice
	hist func(string, ...any)
}
func NewParcel(fn func(string, ...any)) *Parcel {
	var err error
	p := Parcel{ hist: fn }
	if p.voice, err = flite.VoiceSelect("kal"); err != nil {
		p.hist("FAIL flite, %s", err.Error())
	}
	return &p
}
func (p *Parcel) Update() error {

	return nil
}
func (p *Parcel) Experiment() error {
// TODO first babystep is to output WAV file (so we can listen via indep player)
	// flite experiment
	flite.TextToSpeech("Hello World", p.voice, "output.wav")
	return nil
}


