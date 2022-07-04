//go:build !(wasm || js || android)

package acorn

import (
	_ "embed"

	//"log"
)
import (
	"github.com/gen2brain/flite-go"
)

type Acorn struct {
	voice   *flite.Voice
	info func(string, ...any)
}
func NewAcorn(fn func(string, ...any)) *Acorn {
	var err error
	ac := Acorn{ info: fn }
	if ac.voice, err = flite.VoiceSelect("kal"); err != nil {
		ac.info("FAIL flite, %s", err.Error())
	}
	return &ac
}
func (a *Acorn) Update() error {

	/*
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		flite.TextToSpeech("Hello World", g.voice, "play")
	}*/

	return nil
}


