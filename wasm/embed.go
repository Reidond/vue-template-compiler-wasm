package wasm

import (
	_ "embed"
)

//go:embed vue-template-compiler.wasm
var VueTemplateCompilerWasm []byte
