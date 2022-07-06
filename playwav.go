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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	raudio "github.com/hajimehoshi/ebiten/v2/examples/resources/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sampleRate   = 48000
)

type testPlay struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

func NewPlay() (*testPlay, error) {
	g := &testPlay{}

	var err error
	// Initialize audio context.
	g.audioContext = audio.NewContext(sampleRate)

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
	d, err := wav.DecodeWithSampleRate(sampleRate,bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		return nil, err
	}

	// Create an audio.Player that has one stream.
	g.audioPlayer, err = g.audioContext.NewPlayer(d)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *testPlay) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		// As audioPlayer has one stream and remembers the playing position,
		// rewinding is needed before playing when reusing audioPlayer.
		g.audioPlayer.Rewind()
		g.audioPlayer.Play()
	}
	return nil
}


func (g *testPlay) Draw(screen *ebiten.Image) {
	if g.audioPlayer.IsPlaying() {
		ebitenutil.DebugPrint(screen, "Bump!")
	} else {
		ebitenutil.DebugPrint(screen, "Press P to play the wav")
	}
}

