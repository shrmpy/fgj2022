package main

import (
	"image"
	"image/color"
)
import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

func newBurger(wd, ht int, re *etxt.Renderer, fg color.RGBA) *clickable {
	var label = "≡"
	var sz = re.SelectionRect(label)
	return newClickable(wd, 0, etxt.Top, etxt.Right, label, sz, fg)
}

type clickable struct {
	x, y      int
	va        etxt.VertAlign
	ha        etxt.HorzAlign
	Text      string
	mouseDown bool
	rectSize  etxt.RectSize
	onPressed func(el mue)
	fg        color.RGBA
}

func newClickable(x, y int, va etxt.VertAlign, ha etxt.HorzAlign,
	label string, sz etxt.RectSize, fg color.RGBA) *clickable {
	return &clickable{
		x:        x,
		y:        y,
		va:       va,
		ha:       ha,
		Text:     label,
		rectSize: sz,
		fg:       fg,
	}
}

func (c *clickable) Update() {
	// TODO fold in touch events too
	var click = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if !click {
		if c.mouseDown {
			// was depressed, now "key-up"
			c.Action()
		}
		c.mouseDown = false
		return
	}
	var hb = c.HitBox()
	if image.Pt(ebiten.CursorPosition()).In(hb) {
		c.mouseDown = true
	} else {
		c.mouseDown = false
	}
}

// bounds with respect to h/v alignment
func (c *clickable) HitBox() image.Rectangle {
	var minx, miny, maxx, maxy int
	if c.ha == etxt.Left {
		minx = c.x
		maxx = c.x + c.rectSize.WidthCeil()
	} else {
		maxx = c.x
		minx = c.x - c.rectSize.WidthCeil()
	}
	if c.va == etxt.Top {
		miny = c.y
		maxy = c.y + c.rectSize.HeightCeil()
	} else {
		maxy = c.y
		miny = c.y - c.rectSize.HeightCeil()
	}
	return image.Rect(minx, miny, maxx, maxy)
}
func (c *clickable) Draw(re *etxt.Renderer) {
	re.SetAlign(c.va, c.ha)
	re.SetColor(c.fg)
	re.Draw(c.Text, c.x, c.y)
}
func (c *clickable) Action() error {
	if c.onPressed != nil {
		c.onPressed(c)
	}
	return nil
}
func (c *clickable) HandleFunc(f func(el mue)) {
	c.onPressed = f
}

// "minimum UI element" is text that responds to events
type mue interface {
	Action() error
	HandleFunc(f func(el mue))
}
