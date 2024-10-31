package code_engine

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/require"
	"go.uber.org/zap"
)

type Code struct {
	Path    string
	Program *goja.Program
}

type JsVM struct {
	// vm 虚拟机
	vm *goja.Runtime
	// 注册
	registry *require.Registry
	// 事件循环
	loop *EventLoop
	// 全局goja加载目录
	globalPath string
	// 输出回调
	ConsoleCallBack ConsoleCallBack
	// 外部添加到内部的内容
	pkg map[string]map[string]any
	// 对应文件的编译对象
	proMap map[string]*Code
	// ts pro
	tsPro   *goja.Program
	console *console
}

type ModuleFunc = func(vm *goja.Runtime, module *goja.Object)

func NewJsVM(globalPath string, log *zap.Logger, TsPro *goja.Program, consoleCallBack ConsoleCallBack) *JsVM {
	jsVM := &JsVM{
		vm:              goja.New(),
		globalPath:      globalPath,
		ConsoleCallBack: consoleCallBack,
		pkg:             make(map[string]map[string]any),
		proMap:          make(map[string]*Code),
		tsPro:           TsPro,
	}

	// 输出日志
	console := newConsole(log)
	o := jsVM.vm.NewObject()
	o.Set("log", console.Log)
	o.Set("debug", console.Debug)
	o.Set("info", console.Info)
	o.Set("error", console.Error)
	o.Set("warn", console.Warn)
	jsVM.vm.Set("console", o)

	jsVM.console = console
	if consoleCallBack != nil {
		console.CallBack = append(console.CallBack, &consoleCallBack)
	}

	var parserOpts []parser.Option
	parserOpts = append(parserOpts, parser.WithDisableSourceMaps)
	jsVM.vm.SetParserOptions(parserOpts...)

	ops := []require.Option{}

	if globalPath != "" {
		ops = append(ops, require.WithGlobalFolders(globalPath))
	}

	// source
	ops = append(ops, require.WithLoader(jsVM.sourceLoader(globalPath)))

	jsVM.loop = NewEventLoop(jsVM.vm)
	jsVM.registry = require.NewRegistry(ops...)
	jsVM.registry.Enable(jsVM.vm)

	self := jsVM.vm.GlobalObject()
	jsVM.vm.Set("self", self)

	jsVM.vm.Set("atob", func(code string) string {
		raw, err := Base64DecodeString(code)
		if err != nil {
			panic(err)
		}
		return raw
	})

	jsVM.vm.Set("btoa", func(code string) string {
		return Base64EncodeString(code)
	})

	if jsVM.tsPro != nil {
		jsVM.vm.RunProgram(jsVM.tsPro)
	}

	return jsVM
}

func (jv *JsVM) load() {
	// 加载其他模块
	for name, mod := range jv.pkg {
		gojaMod := jv.vm.NewObject()
		for k, v := range mod {
			gojaMod.Set(k, v)
		}

		// 注册模块
		jv.vm.Set(name, gojaMod)
	}
}

// SetProperty
// 向模块写入变量或者写入方法
func (jv *JsVM) SetProperty(moduleName, key string, value any) {
	mod, ok := jv.pkg[moduleName]
	if !ok {
		jv.pkg[moduleName] = make(map[string]any)
		mod = jv.pkg[moduleName]
	}
	mod[key] = value
}

func (jv *JsVM) SetConsoleCallBack(consoleCallBack ConsoleCallBack) {
	jv.console.SetCallBack(&consoleCallBack)
}

func (jv *JsVM) Run(name string, pro *goja.Program) error {

	// 加载模块
	jv.load()

	if pro != nil {
		loop := NewEventLoop(jv.vm)
		var exception error
		loop.Run(func(r *goja.Runtime) {
			_, err := r.RunProgram(pro)
			if gojaErr, ok := err.(*goja.Exception); ok {
				exception = errors.New(gojaErr.String())
				return
			}
		})

		if exception != nil {
			return exception
		}
	} else {
		return errors.New("code is nil")
	}

	return nil
}

func (jv *JsVM) compile(name string, path string) (pro *goja.Program, err error) {
	var tmpPath string
	if jv.globalPath != "" {
		tmpPath = filepath.Join(jv.globalPath, path)
	} else {
		tmpPath = path
	}

	// 读取文件
	var src []byte
	src, err = os.ReadFile(tmpPath)
	if err != nil {
		return nil, err
	}

	code := string(src)

	ext := filepath.Ext(tmpPath)
	if ext == ".ts" {
		src, err = jv.Transpile(src)
		if err != nil {
			return nil, err
		}

		code = string(src)
		code = fmt.Sprintf(`
		function %s() {
				const module = (function(exports) {
					%s
					return exports
				})(self)
				return module.%s();
		}`, name, code, name)
	}

	// fmt.Println("code", code)
	// 编译文件
	pro, err = goja.Compile(name, code, false)
	if err != nil {
		fmt.Println("compile error:", err)
		return nil, err
	}

	jv.proMap[name] = &Code{
		Path:    path,
		Program: pro,
	}

	return
}

func (jv *JsVM) ExportFunc(name string, path string, fn any) error {
	// 编译
	pro, err := jv.compile(name, path)
	if err != nil {
		return err
	}

	vm := jv.vm
	// 运行
	err = jv.Run(name, pro)
	if err != nil {
		return err
	}

	nameFunc := vm.Get(name)

	_, ok := goja.AssertFunction(nameFunc)
	if !ok {
		return fmt.Errorf("%s function not found", name)
	}

	// 导出
	return vm.ExportTo(vm.Get(name), fn)
}
