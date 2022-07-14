
package acorn

import (
	"bytes"
	"log"
)

// TODO for web target, we need to fallback on pre-gen WAV audio?
//      which means UI only gives canned input options for button/click
func FliteSpeech(wasm []byte, speech string) []byte {
	log.Printf("INFO flite enter")
	var b bytes.Buffer

	log.Printf("INFO flite exit")
	return b.Bytes()
}
