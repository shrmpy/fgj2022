
package acorn

import (
	_ "embed"

	//"log"
)

type Acorn struct {
	info func(string, ...any)
}
func NewAcorn(fn func(string, ...any)) *Acorn {
	var err error
	ac := Acorn{ info: fn }
	return &ac
}
func (a *Acorn) Update() error {

	/*
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		flite.TextToSpeech("Hello World", g.voice, "play")
	}*/

	return nil
}


