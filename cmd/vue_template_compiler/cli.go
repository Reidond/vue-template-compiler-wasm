package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Reidond/vue-template-compiler-wasm/wasm"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "fp", "./test.json", "file to read and pass to wasm module")
	flag.Parse()

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Given file does not exist")
	}

	body, err := os.ReadFile(absFilePath)
	if err != nil {
		log.Fatalf("Cannot find given file")
	}

	in := bytes.NewReader(body)
	out := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Combine the above into our baseline config, overriding defaults.
	config := wazero.NewModuleConfig().
		// By default, I/O streams are discarded and there's no file system.
		WithStdin(in).WithStdout(out).WithStderr(stderr)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// InstantiateModule runs the "_start" function, WASI's "main".
	_, err = r.InstantiateWithConfig(ctx, wasm.VueTemplateCompilerWasm, config)
	if err != nil {
		// Note: Most compilers do not exit the module after running "_start",
		// unless there was an error. This allows you to call exported functions.
		if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
			fmt.Fprintf(os.Stderr, "exit_code: %d\n", exitErr.ExitCode())
		} else if !ok {
			log.Panicln(err)
		}
	}

	fmt.Printf("%s", out.String())
}
