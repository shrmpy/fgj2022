package main

import (
	//"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"golang.org/x/image/math/fixed"
)

type cmdbar struct {
	cl      clickable
	infocus bool
	bkspc   bool
	quit    func(mue)
	tts     func(string)
	input   []rune
}

func newBar(wd, ht int, re *etxt.Renderer, fg color.RGBA) *cmdbar {
	var label = ">"
	var sz = re.SelectionRect(label)
	// stretch width to full row
	sz.Width = fixed.I(wd)
	var bar = &cmdbar{input: []rune{}}
	var cl = newClickable(0, ht, etxt.Bottom, etxt.Left, label, sz, fg)
	cl.HandleFunc(bar.dispatch)
	bar.cl = *cl
	return bar
}

func (b *cmdbar) Update(touches []ebiten.TouchID) error {
	if b.infocus && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// accept text input when enter is pressed
		b.infocus = false
		b.cl.Action()
		return nil
	}
	if b.cl.hasTouch(touches) {
		b.infocus = true
	}
	if b.infocus {
		var del = ebiten.IsKeyPressed(ebiten.KeyBackspace)
		if del && !b.bkspc {
			b.bkspc = true
			return nil
		} else if !del {
			if b.bkspc == true {
				//key-up
				b.bkspc = false
				var sz = len(b.input) - 1
				if sz == 0 {
					b.input = []rune{}
				} else if sz > 0 {
					b.input = b.input[:sz-1]
				}
				return nil
			}
			b.bkspc = false
		}
		b.input = ebiten.AppendInputChars(b.input)
	}
	return nil
}

func (b *cmdbar) Draw(re *etxt.Renderer) {
	b.cl.Text = ">" + string(b.input)
	b.cl.Draw(re)
}

func (b *cmdbar) QuitFunc(fn func(el mue)) {
	b.quit = fn
}
func (b *cmdbar) TTSFunc(fn func(sp string)) {
	b.tts = fn
}
func (b *cmdbar) dispatch(el mue) {
	// TODO extract last text input line
	//TODO pass input to flite convert
	var op = strings.ToLower(strings.TrimPrefix(b.cl.Text, ">"))
	if op == "/quit" {
		if b.quit != nil {
			b.quit(el)
		}
		return
	}
	if b.tts != nil {
		b.tts(b.cl.Text)
	}
}
