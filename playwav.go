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
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	//"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	//raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	riaudio "github.com/hajimehoshi/ebiten/v2/examples/resources/images/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

const (
	sampleRate   = 16000
)
var alertButtonImage *ebiten.Image

type testPlay struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	arrow *clickable
	game *Game
	alertButtonPosition image.Point
}

func NewPlay(g *Game, wd, ht int, re *etxt.Renderer) (*testPlay, error) {
	var w = &testPlay{ game: g }

	var err error
	// Initialize audio context.
	w.audioContext = audio.NewContext(sampleRate)

	// In this example, embedded resource "Jab_wav" is used.
	//
	// If you want to use a wav file, open this and pass the file stream to wav.Decode.
	// Note that file's Close() should not be closed here
	// since audio.Player manages stream state.
	//
	//     f, err := os.Open("jab.wav")
	//     if err != nil {
	//         return err
	//     }
	//
	//     d, err := wav.DecodeWithoutResampling(f)
	//     ...

	// Decode wav-formatted data and retrieve decoded PCM stream.
	d, err := wav.DecodeWithSampleRate(sampleRate,bytes.NewReader(flite_WAV))
	//d, err := vorbis.DecodeWithSampleRate(sampleRate,bytes.NewReader(raudio.Ragtime_ogg))
	if err != nil {
		return nil, err
	}

	// Create an audio.Player that has one stream.
	w.audioPlayer, err = w.audioContext.NewPlayer(d)
	if err != nil {
		return nil, err
	}

	w.arrow = newArrow(wd,ht,re,color.RGBA{0xff,0x72,0x5c,0xff})
	w.arrow.HandleFunc(func(el mue) { w.toggleAudio() })

	w.audioPlayer.SetVolume(1)
	const btnPadding = 16

	img, _, err := image.Decode(bytes.NewReader(riaudio.Alert_png))
	if err != nil {
		return nil,err
	}
	alertButtonImage = ebiten.NewImageFromImage(img)
	var sz, _ = alertButtonImage.Size()

	w.alertButtonPosition.X = (ht -sz*2 + btnPadding*1) / 2 + sz + btnPadding
	w.alertButtonPosition.Y = ht - 160

	return w, nil
}

func newArrow(wd, ht int, re *etxt.Renderer, fg color.RGBA) *clickable {
        var label = "▶PLAY"
        var sz = re.SelectionRect(label)
        return newClickable(0, ht, etxt.Bottom, etxt.Left, label, sz, fg)
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
		//ctx := context.WithTimeout(context.Background())
		reader, writer := io.Pipe()
		log.Printf("INFO flite enter")
		go func(){
			fliteTest(writer, "Flite hello world placeholder.")
			log.Printf("INFO flite exit")
		}()

	        ////if _, err := wav.DecodeWithSampleRate(sampleRate,data); err != nil {
	        srd, err := wav.DecodeWithSampleRate(sampleRate,reader)
		if err != nil {
			log.Printf("DEBUG buffer, %s", err.Error())
		}
		log.Printf("INFO decode, %d", srd.Length())
		// TODO signal goroutine (context-with-timeout for safety?)
		//ctx.Cancel()
		writer.Close()
	}

	return nil
}

func (w *testPlay) Draw(re *etxt.Renderer, screen *ebiten.Image) {

	if w.audioPlayer.IsPlaying() {
		log.Printf("INFO playing")
		w.arrow.Text ="■STOP"
	} else {
		w.arrow.Text ="▶PLAY"
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
func (w *testPlay) Close() {
	if w.audioPlayer != nil {
		w.audioPlayer.Close()
	}
}

