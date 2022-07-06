//go:build !(wasm || js || android)

package acorn

import (

	//"log"
)

type Parcel struct {
	hist func(string, ...any)
}
func NewParcel(fn func(string, ...any)) *Parcel {
	p := Parcel{ hist: fn }
	return &p
}
func (p *Parcel) Update() error {

	return nil
}
func (p *Parcel) Experiment() error {
// TODO first babystep is to output WAV file (so we can listen via indep player)
	////flite.TextToSpeech("Hello World", p.voice, "output.wav")
	return nil
}


