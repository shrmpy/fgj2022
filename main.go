package main

import (
	_ "embed"

	"fmt"
	"image/color"
	"log"
)
import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"github.com/shrmpy/fgj2022/acorn"
)

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js dist/web/wasm_exec.js
//go:generate env GOOS=js GOARCH=wasm go build -ldflags "-w -s" -o dist/web/fgj2022.wasm ./
//go:embed DejaVuSansMono.ttf
var dejavuSansMonoTTF []byte

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}
func main() {
	var (
		err    error
		wd, ht = 640, 480
	)
	var game = &Game{
		Width:   wd,
		Height:  ht,
		txtre:   newRenderer(),
		history: make([]string, 0, 25),
	}
	game.p = acorn.NewParcel(game.AddHistory)
	if game.play, err = NewPlay(wd,ht,game.txtre); err != nil {
		log.Fatalf("FAIL wav, %s", err.Error())
	}

	ebiten.SetWindowSize(wd, ht)
	ebiten.SetWindowTitle("fgj2022")
	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("FAIL main, %s", err.Error())
	}
}

// Update runs game logic steps
func (g *Game) Update() error {
	// Pressing Q any time quits immediately
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return fmt.Errorf("INFO Quit key")
	}

	// Pressing F toggles full-screen
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		var fs = ebiten.IsFullscreen()
		ebiten.SetFullscreen(!fs)
	}

	// TODO 
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.p.Experiment()
	}
	g.play.Update()

	return nil
}

// Draw renders one frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.printHistory(screen)
	g.play.Draw(g.txtre)

}

func newRenderer() *etxt.Renderer {
	var (
		err error
		name string
		fonts  = etxt.NewFontLibrary()
	)
	if name, err = fonts.ParseFontBytes(dejavuSansMonoTTF); err != nil {
		log.Fatalf("FAIL Parse DejaVu Sans Mono, %s", err.Error())
	}
	log.Printf("INFO font, %s", name)
	var renderer = etxt.NewStdRenderer()
	renderer.SetCacheHandler(etxt.NewDefaultCache(2 * 1024 * 1024).NewHandler())
	renderer.SetFont(fonts.GetFont("DejaVu Sans Mono"))
	renderer.SetColor(color.White)
	renderer.SetSizePx(18)
	return renderer
}
func (g *Game) printHistory(screen *ebiten.Image) {
	// help us troubleshoot
	g.txtre.SetTarget(screen)
	max := len(g.history)

	g.txtre.SetAlign(etxt.Bottom, etxt.Left)
	for i := max; i > 0; i-- {
		msg := g.history[i-1]
		sz := g.txtre.SelectionRect(msg)
		g.txtre.Draw(msg, 0, g.Height-sz.Height.Ceil()*i)
	}
	// print frame rate in se corner
	g.txtre.SetAlign(etxt.Bottom, etxt.Right)
	g.txtre.Draw(fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), g.Width-1, g.Height)
}

// Game represents the main game state
type Game struct {
	Width  int
	Height int

	txtre   *etxt.Renderer
	history []string
	p   *acorn.Parcel
	play *testPlay
}

// Layout is static for now, can be dynamic
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

// allow maze to bubble-up debug msg
func (g *Game) AddHistory(tmp string, values ...any) {
	var msg = fmt.Sprintf(tmp, values...)
	if len(g.history) < 25 {
		g.history = append(g.history, msg)
	}
}
