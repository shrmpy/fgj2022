
package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"golang.org/x/image/math/fixed"
)

type cmdbar struct {
	clickable
	infocus bool
}

func newBar(wd, ht int, re *etxt.Renderer, fg color.RGBA) *cmdbar {
	var label = ">"
	var sz = re.SelectionRect(label)
	// stretch width to full row
	sz.Width = fixed.I(wd)
	var cl = newClickable(0, ht, etxt.Bottom, etxt.Left, label, sz, fg)
	return &cmdbar{clickable: *cl}
}

func (b *cmdbar) Update(touches []ebiten.TouchID) error {
	if b.infocus && inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// accept text input when enter is pressed
		b.infocus = false
		//TODO pass input to flite convert
		return nil
	}
	if b.hasFocus(touches) {
		b.infocus = true
	}
	if b.infocus {
		var rs = ebiten.AppendInputChars([]rune(b.clickable.Text))
		b.clickable.Text = string(rs)
	}
	return nil
}

func (b *cmdbar) Draw(re *etxt.Renderer) {
	b.clickable.Draw(re)
}
// any touch event?
func (b *cmdbar) hasFocus(touches []ebiten.TouchID) bool {
        var r = b.clickable.HitBox()
        if image.Pt(ebiten.CursorPosition()).In(r) {
                if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
                        return true
                }
        }

        for _, id := range touches {
                if image.Pt(ebiten.TouchPosition(id)).In(r) {
                        return true
                }
        }
        return false
}


