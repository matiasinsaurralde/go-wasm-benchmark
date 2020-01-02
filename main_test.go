package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	wagon_wasm "github.com/go-interpreter/wagon/wasm"
	life_wasm "github.com/perlin-network/life/exec"
	wasmer_wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

const (
	funcName = "sum"
)

var (
	simpleWasmBytes []byte
)

func init() {
	simpleWasmBytes, _ = ioutil.ReadFile("simple.wasm")
}

func BenchmarkWasmerSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		instance, err := wasmer_wasm.NewInstance(simpleWasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		defer instance.Close()
		sum := instance.Exports[funcName]
		_, err = sum(5, 37)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWasmerSumReentrant(b *testing.B) {
	instance, err := wasmer_wasm.NewInstance(simpleWasmBytes)
	if err != nil {
		panic(err)
	}
	defer instance.Close()
	sum := instance.Exports[funcName]
	for n := 0; n < b.N; n++ {
		_, err := sum(5, 37)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkWagonSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		wagon_wasm.SetDebugMode(false)
		reader := bytes.NewReader(simpleWasmBytes)
		m, err := wagon_wasm.ReadModule(reader, wagonImporter)
		if err != nil {
			b.Fatal(err)
		}
		if m.Export == nil {
			b.Fatalf("module has no export section")
		}
		vm, err := exec.NewVM(m)
		if err != nil {
			b.Fatalf("could not create VM: %v", err)
		}
		defer vm.Close()

		var funcCode int64
		for name, e := range m.Export.Entries {
			funcCode = int64(e.Index)
			if name == funcName {
				break
			}
		}
		vm.ExecCode(funcCode, 5, 37)
	}
}

func BenchmarkWagonSumReentrant(b *testing.B) {
	wagon_wasm.SetDebugMode(false)
	reader := bytes.NewReader(simpleWasmBytes)
	m, err := wagon_wasm.ReadModule(reader, wagonImporter)
	if err != nil {
		panic(err)
	}
	if m.Export == nil {
		log.Fatalf("module has no export section")
	}
	vm, err := exec.NewVM(m)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}
	defer vm.Close()
	var funcCode int64
	for name, e := range m.Export.Entries {
		funcCode = int64(e.Index)
		if name == funcName {
			break
		}
	}
	for n := 0; n < b.N; n++ {
		vm.ExecCode(funcCode, 5, 37)
		vm.Restart()
	}
}

func wagonImporter(name string) (*wagon_wasm.Module, error) {
	return nil, nil
}

func BenchmarkLifeSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		vm, err := life_wasm.NewVirtualMachine(simpleWasmBytes, life_wasm.VMConfig{}, &life_wasm.NopResolver{}, nil)
		if err != nil {
			b.Fatal(err)
		}
		entryID, ok := vm.GetFunctionExport(funcName)
		if !ok {
			b.Fatal(err)
		}
		_, err = vm.Run(entryID, 5, 37)
		if err != nil {
			b.Fatal(err)
		}
	}
}
