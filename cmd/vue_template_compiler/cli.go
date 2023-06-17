package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Reidond/vue-template-compiler-wasm/internal/compiler"
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

	var wasmInput compiler.WasmInput
	err = json.Unmarshal(body, &wasmInput)
	if err != nil {
		log.Fatalf("Cannot unmarshal given file, Err: %v", err)
	}

	out, err := compiler.CompileSfcCode(wasmInput.SfcCode, wasmInput.MountID)
	if err != nil {
		log.Fatalf("Cannot compile given file, Err: %v", err)
	}

	outMarshal, err := json.Marshal(out)
	if err != nil {
		log.Fatalf("Cannot marshal given file, Err: %v", err)
	}

	fmt.Printf("%s", string(outMarshal))
}
