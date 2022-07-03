package main

import (
	_ "embed"

	"fmt"
	"image/color"
	"log"
)
import (
	////"github.com/gen2brain/flite-go"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
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
		name   string
		wd, ht = 640, 480
		ch     = make(chan string, 100)
		fonts  = etxt.NewFontLibrary()
	)
	defer close(ch)
	if name, err = fonts.ParseFontBytes(dejavuSansMonoTTF); err != nil {
		log.Fatalf("FAIL Parse DejaVu Sans Mono, %s", err.Error())
	}
	log.Printf("INFO font, %s", name)
	var renderer = etxt.NewStdRenderer()
	renderer.SetCacheHandler(etxt.NewDefaultCache(2 * 1024 * 1024).NewHandler())
	renderer.SetFont(fonts.GetFont("DejaVu Sans Mono"))
	renderer.SetColor(color.White)
	renderer.SetSizePx(12)
	ebiten.SetWindowSize(wd, ht)
	ebiten.SetWindowTitle("fgj2022")
	var game = &Game{
		Width:   wd,
		Height:  ht,
		txtre:   renderer,
		history: make([]string, 0, 25),
	}/*
	if game.voice, err = flite.VoiceSelect("kal"); err != nil {
		log.Fatalf("FAIL flite, %s", err.Error())
	}*/

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

	// TODO wasi for web target
	/*
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		flite.TextToSpeech("Hello World", g.voice, "play")
	}*/

	return nil
}

// Draw renders one frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.printHistory(screen)

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
	////voice   *flite.Voice
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
