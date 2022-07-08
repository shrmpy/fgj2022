//go:build !(s390x && mips)

package main

import (
	"context"
	"embed"
	_ "embed"
	"log"
	"io/fs"
	"os"
)
import (
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)
//go:embed dist/web
var flite_fs embed.FS
//go:embed dist/testout.wav
var flite_WAV []byte
//go:embed dist/flite.wasm
var flite_wasm []byte

func fliteTest(speak string) {
	var ctx = context.Background()
	var rt = wazero.NewRuntimeWithConfig(wazero.NewRuntimeConfig().WithWasmCore2())
	defer rt.Close(ctx)
	tmp, err := fs.Sub(flite_fs, "dist")
	if err != nil {
		log.Panicf("FAIL fs, %s", err.Error())
	}
	var cfg = wazero.NewModuleConfig().
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		WithFS(tmp)

	if _, err := wasi_snapshot_preview1.Instantiate(ctx, rt); err != nil {
		log.Panicf("FAIL wazero, %s", err.Error())
	}
	code, err := rt.CompileModule(ctx, flite_wasm, wazero.NewCompileConfig())
	if err != nil {
		log.Panicf("FAIL compile, %s", err.Error())
	}

	// InstantiateModule runs the "_start" function, WASI's "main".
	// * Set the program name (arg[0]) to "wasi"; arg[1] should be "/test.txt".
	cfg.WithArgs("wasi", speak, "output.wav")
	if _, err = rt.InstantiateModule(ctx, code, cfg); err != nil {

		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an error. This allows you to call exported functions.
		if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
			log.Printf("INFO exit_code: %d", exitErr.ExitCode())
		} else if !ok {
			log.Panicf("FAIL wasi, %s", err.Error())
		}
	}

/*	fli, err := rt.InstantiateModule(ctx, flite_wasm, cfg)
	results, err := xfn.Call(ctx, args, "output.wav")
	*/
}

