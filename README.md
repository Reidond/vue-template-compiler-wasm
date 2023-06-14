# vue-template-compiler-wasm

This is experiment to get `vue-template-compiler` from Vue 2 toolchain to work as wasm module through compiling js bundle using Bytecode Alliance Javy wasm compiler.

After that was done, I tried to run wasm module in go, using wazero runtime was only option because currently wasmtime doesn't have bindings to set stdin to some buffer. 
