package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/go-interpreter/wagon/exec"
	wagon_wasm "github.com/go-interpreter/wagon/wasm"
	wasm3 "github.com/matiasinsaurralde/go-wasm3"
	life_wasm "github.com/perlin-network/life/exec"
	wasmer_wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

const (
	funcName        = "sum"
	cstringFuncName = "somecall"
)

var (
	simpleWasmBytes  []byte
	cstringWasmBytes []byte
)

func init() {
	simpleWasmBytes, _ = ioutil.ReadFile("simple.wasm")
	cstringWasmBytes, _ = ioutil.ReadFile("cstring.wasm")
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

func BenchmarkWASM3Sum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		env := wasm3.NewEnvironment()
		defer env.Destroy()
		runtime := wasm3.NewRuntime(env, 64*1024)
		defer runtime.Destroy()
		_, err := runtime.Load(simpleWasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		fn, err := runtime.FindFunction(funcName)
		if err != nil {
			b.Fatal(err)
		}
		fn(5, 37)
	}
}

func BenchmarkWASM3SumReentrant(b *testing.B) {
	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
	defer runtime.Destroy()
	_, err := runtime.Load(simpleWasmBytes)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		fn, err := runtime.FindFunction(funcName)
		if err != nil {
			b.Fatal(err)
		}
		fn(5, 37)
	}
}

func BenchmarkWASM3CString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		env := wasm3.NewEnvironment()
		defer env.Destroy()
		runtime := wasm3.NewRuntime(env, 64*1024)
		defer runtime.Destroy()
		_, err := runtime.Load(cstringWasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		fn, err := runtime.FindFunction(cstringFuncName)
		if err != nil {
			b.Fatal(err)
		}
		result := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()
		mem := runtime.GetMemory(memoryLength, 0)
		buf := new(bytes.Buffer)
		for n := 0; n < memoryLength; n++ {
			if n < result {
				continue
			}
			value := mem[n]
			if value == 0 {
				break
			}
			buf.WriteByte(value)
		}
	}
}

func BenchmarkWASM3CStringReentrant(b *testing.B) {
	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
	defer runtime.Destroy()
	_, err := runtime.Load(cstringWasmBytes)
	if err != nil {
		b.Fatal(err)
	}
	fn, err := runtime.FindFunction(cstringFuncName)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		result := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()
		mem := runtime.GetMemory(memoryLength, 0)
		buf := new(bytes.Buffer)
		for n := 0; n < memoryLength; n++ {
			if n < result {
				continue
			}
			value := mem[n]
			if value == 0 {
				break
			}
			buf.WriteByte(value)
		}
	}
}

func BenchmarkWasmerCString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		instance, err := wasmer_wasm.NewInstance(cstringWasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		defer instance.Close()
		somecall := instance.Exports[cstringFuncName]
		ptr, err := somecall()
		if err != nil {
			b.Fatal(err)
		}
		mem := instance.Memory.Data()
		buf := new(bytes.Buffer)
		i := int(ptr.ToI32())
		for n := 0; n < len(mem); n++ {
			if n < i {
				continue
			}
			value := mem[n]
			if value == 0 {
				break
			}
			buf.WriteByte(value)
		}
	}
}
