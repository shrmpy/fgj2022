//go:build !(s390x || mips || js)

package acorn

import (
	"bytes"
	"context"
	_ "embed"
	"log"
	"os"
)
import (
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func FliteSpeech(wasm []byte, speech string) []byte {
	log.Printf("INFO flite enter")

	var ctx = context.Background()
	var rt = wazero.NewRuntime(ctx)
	defer rt.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, rt)
	code, err := rt.CompileModule(ctx, wasm)
	if err != nil {
		log.Panicf("FAIL compile, %s", err.Error())
	}
	var b bytes.Buffer
	var cfg = wazero.NewModuleConfig().
		WithStdout(&b).
		WithStderr(os.Stderr)
	// InstantiateModule runs the "_start" function, WASI's "main".
	// * Set the program name (arg[0]) to "wasi"; arg[1] should be "/test.txt".
	cfg.WithArgs("wasi", speech)
	if _, err := rt.InstantiateModule(ctx, code, cfg); err != nil {
		log.Panicf("FAIL wasi, %s", err.Error())
	}
	log.Printf("INFO flite exit")
	return b.Bytes()
}
