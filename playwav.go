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
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	//"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/tinne26/etxt"
)

const (
	sampleRate   = 32000
)

type testPlay struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
	arrow *clickable
	game *Game
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
	//d, err := wav.DecodeWithSampleRate(sampleRate,bytes.NewReader(raudio.Jab_wav))
	d, err := vorbis.DecodeWithSampleRate(sampleRate,bytes.NewReader(raudio.Ragtime_ogg))
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

	return w, nil
}

func newArrow(wd, ht int, re *etxt.Renderer, fg color.RGBA) *clickable {
        var label = "▶"
        var sz = re.SelectionRect(label)
        return newClickable(0, ht, etxt.Bottom, etxt.Left, label, sz, fg)
}

func (w *testPlay) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		// As audioPlayer has one stream and remembers the playing position,
		// rewinding is needed before playing when reusing audioPlayer.
		//w.audioPlayer.Rewind()
		//w.audioPlayer.Play()
		w.arrow.Action()
		return nil
	}
	//w.arrow.Update()
	// search touch events
	if w.controlAudio() {
		w.arrow.Action()
	}

	return nil
}

func (w *testPlay) Draw(re *etxt.Renderer) {
	if w.audioPlayer.IsPlaying() {
		log.Printf("INFO playing")
		w.arrow.Text ="■"
	} else {
		w.arrow.Text ="▶"
	}
	//ebitenutil.DebugPrint(screen, "Press P to play the wav")
	w.arrow.Draw(re)
}
// any touch event?
func (w *testPlay) controlAudio() bool {
/*
	r := image.Rectangle{
		Min: p.playButtonPosition,
		Max: p.playButtonPosition.Add(image.Pt(playButtonImage.Size())),
	}*/
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
// attached to the play-arrow
func (w *testPlay) toggleAudio() {
	if w.audioPlayer.IsPlaying() {
		w.audioPlayer.Pause()
		return
	}
	w.audioPlayer.Play()
}
func (w *testPlay) Close() {
	if w.audioPlayer != nil {
		w.audioPlayer.Close()
	}
}

