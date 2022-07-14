// Copyright 2017 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	riaudio "github.com/hajimehoshi/ebiten/v2/examples/resources/images/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"github.com/shrmpy/fgj2022/acorn"
)

//go:embed dist/testout.wav
var flite_WAV []byte

//go:embed dist/flite.wasm
var flite_wasm []byte

const (
	sampleRate = 16000
)

var alertButtonImage *ebiten.Image

type testPlay struct {
	audioContext        *audio.Context
	audioPlayer         *audio.Player
	arrow               *clickable
	game                *Game
	alertButtonPosition image.Point
}

func NewPlay(g *Game, wd, ht int, re *etxt.Renderer) (*testPlay, error) {
	var w = &testPlay{game: g}
	var err error
	// Initialize audio context.
	w.audioContext = audio.NewContext(sampleRate)
	// Decode wav-formatted data and retrieve decoded PCM stream.
	d, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(flite_WAV))
	if err != nil {
		return nil, err
	}
	if err = w.replaceSource(d); err != nil {
		return nil, err
	}

	w.arrow = newArrow(wd, ht, re, color.RGBA{0xff, 0x72, 0x5c, 0xff})
	w.arrow.HandleFunc(func(el mue) { w.toggleAudio() })

	const btnPadding = 16

	img, _, err := image.Decode(bytes.NewReader(riaudio.Alert_png))
	if err != nil {
		return nil, err
	}
	alertButtonImage = ebiten.NewImageFromImage(img)
	var sz, _ = alertButtonImage.Size()

	w.alertButtonPosition.X = (ht-sz*2+btnPadding*1)/2 + sz + btnPadding
	w.alertButtonPosition.Y = ht - 160

	return w, nil
}

func newArrow(wd, ht int, re *etxt.Renderer, fg color.RGBA) *clickable {
	var label = "▶PLAY"
	var sz = re.SelectionRect(label)
	var aboveBottom = ht - sz.HeightCeil()
	return newClickable(0, aboveBottom, etxt.Bottom, etxt.Left, label, sz, fg)
}

func (w *testPlay) Update() error {
	/* //TODO press P key needs keyboard in mobile
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		w.arrow.Action()
		return nil
	}*/

	// search touch events
	if w.controlAudio() {
		w.arrow.Action()
	}
	if w.fliteAudio() {
		var buf = acorn.FliteSpeech(flite_wasm, "Flite hello world placeholder.")
		str, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(buf))
		if err != nil {
			log.Printf("DEBUG buffer, %s", err.Error())
		}
		log.Printf("INFO decode, %d", str.Length())
		if err := w.replaceSource(str); err != nil {
			log.Printf("DEBUG source, %s", err.Error())
		}
	}

	return nil
}

func (w *testPlay) Draw(re *etxt.Renderer, screen *ebiten.Image) {

	if w.audioPlayer.IsPlaying() {
		log.Printf("INFO playing")
		w.arrow.Text = "■STOP"
	} else {
		w.arrow.Text = "▶PLAY"
	}

	w.arrow.Draw(re)

	var op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(w.alertButtonPosition.X), float64(w.alertButtonPosition.Y))
	screen.DrawImage(alertButtonImage, op)
}

// any touch event?
func (w *testPlay) controlAudio() bool {
	var r = w.arrow.HitBox()
	if image.Pt(ebiten.CursorPosition()).In(r) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	//todo alt instead of ref to game
	for _, id := range w.game.justPressedTouchIDs {
		if image.Pt(ebiten.TouchPosition(id)).In(r) {
			return true
		}
	}
	return false
}
func (w *testPlay) fliteAudio() bool {
	var r = image.Rectangle{
		Min: w.alertButtonPosition,
		Max: w.alertButtonPosition.Add(image.Pt(alertButtonImage.Size())),
	}

	if image.Pt(ebiten.CursorPosition()).In(r) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	//todo alt instead of ref to game
	for _, id := range w.game.justPressedTouchIDs {
		if image.Pt(ebiten.TouchPosition(id)).In(r) {
			return true
		}
	}
	return false
}

// attached to the play-arrow
func (w *testPlay) toggleAudio() {
	if w.audioPlayer.IsPlaying() {
		w.audioPlayer.Pause()
		return
	}
	w.audioPlayer.Rewind()
	w.audioPlayer.Play()
}
func (w *testPlay) replaceSource(str *wav.Stream) error {
	w.Close()
	var err error
	// Create an audio.Player that has one stream.
	w.audioPlayer, err = w.audioContext.NewPlayer(str)
	if err != nil {
		return err
	}
	w.audioPlayer.SetVolume(1)
	return nil
}
func (w *testPlay) Close() {
	if w.audioPlayer != nil {
		w.audioPlayer.Close()
	}
}
