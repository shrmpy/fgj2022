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
	"github.com/shrmpy/fgj2022/acorn"
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
		wd, ht = 640, 480
		ch     = make(chan signal, 100)
	)
	defer close(ch)
	var game = &Game{
		Width:   wd,
		Height:  ht,
		txtre:   newRenderer(),
		history: make([]string, 0, 25),
		bus:     ch,
	}
	game.p = acorn.NewParcel(game.AddHistory)
	if game.play, err = NewPlay(game, wd, ht, game.txtre); err != nil {
		log.Fatalf("FAIL wav, %s", err.Error())
	}
	defer game.play.Close()
	game.bar = newBar(wd, ht, game.txtre, color.RGBA{0x66, 0x33, 0x99, 0xff})
	game.bar.QuitFunc(game.quitSignal)
	game.bar.TTSFunc(game.speechText)
	game.burger = newBurger(wd, ht, game.txtre, color.RGBA{0x66, 0x33, 0x99, 0xff})
	game.burger.HandleFunc(game.quitSignal)

	ebiten.SetWindowSize(wd, ht)
	ebiten.SetWindowTitle("fgj2022")
	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("FAIL main, %s", err.Error())
	}
}

// Update runs game logic steps
func (g *Game) Update() error {
	/*
		// Pressing F toggles full-screen
		if inpututil.IsKeyJustPressed(ebiten.KeyF) {
			var fs = ebiten.IsFullscreen()
			ebiten.SetFullscreen(!fs)
		}*/

	select {
	case req := <-g.bus:
		if req.op == 8888 {
			return fmt.Errorf("INFO Teardown")
		}

	default:
		/*// TODO
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.p.Experiment()
		}*/
		g.justPressedTouchIDs = inpututil.AppendJustPressedTouchIDs(g.justPressedTouchIDs[:0])

		g.burger.Update(g.justPressedTouchIDs)
		g.play.Update(g.justPressedTouchIDs)
		g.bar.Update(g.justPressedTouchIDs)

	}

	return nil
}

// Draw renders one frame
func (g *Game) Draw(screen *ebiten.Image) {
	g.printHistory(screen)
	g.play.Draw(g.txtre, screen)
	g.bar.Draw(g.txtre)
	g.burger.Draw(g.txtre)
}

func newRenderer() *etxt.Renderer {
	var (
		err   error
		name  string
		fonts = etxt.NewFontLibrary()
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
	/*
		// print frame rate in se corner
		g.txtre.SetAlign(etxt.Bottom, etxt.Right)
		g.txtre.Draw(fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()), g.Width-1, g.Height)
	*/
}

// Game represents the main game state
type Game struct {
	Width  int
	Height int

	txtre               *etxt.Renderer
	history             []string
	p                   *acorn.Parcel
	play                *testPlay
	justPressedTouchIDs []ebiten.TouchID
	bar                 *cmdbar
	bus                 chan signal
	burger              *clickable
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

func (g *Game) quitSignal(el mue) {
	log.Printf("INFO SIG quit")
	g.bus <- signal{op: 8888}
}
func (g *Game) speechText(sp string) {
	g.play.Speech(sp)
}

type signal struct {
	op   int
	data string
}
