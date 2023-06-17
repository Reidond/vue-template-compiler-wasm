package compiler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Reidond/vue-template-compiler-wasm/wasm"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/tetratelabs/wazero/sys"
)

type WasmInput struct {
	MountID string `json:"mountId"`
	SfcCode string `json:"sfcCode"`
}

type WasmOutput struct {
	AppCode   string   `json:"appCode"`
	StyleCode []string `json:"styleCode"`
}

func CompileSfcCode(sfcCode string, mountID string) (WasmOutput, error) {
	empty := WasmOutput{}

	input := WasmInput{
		MountID: mountID,
		SfcCode: sfcCode,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return empty, err
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
			return empty, fmt.Errorf("exit_code: %d", exitErr.ExitCode())
		} else if !ok {
			return empty, err
		}
	}

	var wasmOutput WasmOutput
	err = json.Unmarshal(out.Bytes(), &wasmOutput)
	if err != nil {
		return empty, err
	}

	return wasmOutput, nil
}
