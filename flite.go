//go:build !(s390x && mips)

package main

import (
	"bytes"
	"context"
	_ "embed"
	"log"
	//"io/fs"
	"os"
)
import (
	"github.com/tetratelabs/wazero"
	//"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)


//go:embed dist/testout.wav
var flite_WAV []byte
//go:embed dist/flite.wasm
var flite_wasm []byte

func fliteTest(speech string) *bytes.Buffer {
	var ctx = context.Background()
	var rt = wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2())
	defer rt.Close(ctx)
	var data bytes.Buffer
	var cfg = wazero.NewModuleConfig().
		WithStdout(&data).
		WithStderr(os.Stderr)


	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		log.Panicf("FAIL wazero, %s", err.Error())
	}

	code, err := rt.CompileModule(ctx, flite_wasm, wazero.NewCompileConfig())
	if err != nil {
		log.Panicf("FAIL compile, %s", err.Error())
	}

	cfg.WithArgs("wasi", speech)
	_, err = rt.InstantiateModule(ctx, code, cfg)
	if err != nil {
		log.Panicf("FAIL wasi, %s", err.Error())
	}

	// InstantiateModule runs the "_start" function, WASI's "main".
	// * Set the program name (arg[0]) to "wasi"; arg[1] should be "/test.txt".

	return &data
}


