package code_engine

import (
	_ "embed"
)

// ts v4.9.3.js

//go:embed ts.js
var tsTranspile []byte

func (vm *JsVM) TranspileFile(file string) ([]byte, error) {
	code, err := ReadFile(file)
	if err != nil {
		return nil, err
	}

	return vm.Transpile(code)
}

func (vm *JsVM) Transpile(code []byte) ([]byte, error) {
	compiler := `{"strict":false,"target":"ES5","module":"CommonJS"}`

	s := `ts.transpileModule(atob("` + Base64EncodeString(Bytes2String(code)) + `"), {
		"compilerOptions": ` + compiler + `,
		})`

	res, err := vm.vm.RunString(s)
	if err != nil {
		return nil, err
	}

	return String2Bytes(ToString(res.ToObject(vm.vm).Get("outputText").Export()) + ""), nil
}
