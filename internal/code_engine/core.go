package code_engine

import (
	"sync"

	goja "github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja_nodejs/require"
	"go.uber.org/zap"
)

// Core
// goja运行时核心的结构体
type Core struct {
	// 对象池
	pool sync.Pool
	// 全局goja加载目录
	globalPath string
	// 外部注册的模块
	modules map[string]ModuleFunc
	// 日志对象
	logger *zap.Logger
	// 日志存放路径
	logPath string
	// 日志输出模式
	logMode LogOutMode
	// 输出回调
	ConsoleCallBack ConsoleCallBack
}

type OptFunc = func(*Core)

func NewCore(opts ...OptFunc) *Core {
	c := new(Core)
	c.modules = make(map[string]ModuleFunc)

	for _, opt := range opts {
		opt(c)
	}

	// 注册原生模块
	for Name, moduleFn := range c.modules {
		require.RegisterNativeModule(Name, func(runtime *goja.Runtime, module *goja.Object) {
			m := module.Get("exports").(*goja.Object)
			moduleFn(runtime, m)
		})
	}

	c.InitPool()

	return c
}

func (c *Core) InitPool() {
	pro, err := c.CompileTS()
	if err == nil {
		c.pool = sync.Pool{
			New: func() interface{} {
				jsVM := NewJsVM(
					c.globalPath,
					c.logger,
					pro,
					c.ConsoleCallBack,
				)

				return jsVM
			},
		}
	}
}

// OptionLog
// 配置日志
func OptionLog(path string, log *zap.Logger) OptFunc {
	return func(core *Core) {
		core.logger = log
		core.logPath = path
	}
}

func (c *Core) GetRuntime() *JsVM {
	vm := c.pool.Get().(*JsVM)
	return vm
}

func (c *Core) PutRuntime(vm *JsVM) {
	c.pool.Put(vm)
}

func (c *Core) CompileTS() (prg *goja.Program, err error) {
	var astPrg *ast.Program
	astPrg, err = goja.Parse("", Bytes2String(tsTranspile), parser.WithDisableSourceMaps)
	if err != nil {
		return
	}

	prg, err = goja.CompileAST(astPrg, true)
	return
}

// SetLogPath
// 设置日志路径
func (c *Core) SetLogPath(path string) {
	c.logPath = path
}

// SetLogOutMode
// 日志输出模式
// debug 输出到控制台和输出到日志文件
// release 只输出到日志文件
func (c *Core) SetLogOutMode(mod LogOutMode) {
	c.logMode = mod
}

func (c *Core) SetGlobalPath(path string) {
	c.globalPath = path
}

// RegisterModule
// 注册模块
func (c *Core) RegisterModule(moduleName string, fn ModuleFunc) {
	c.modules[moduleName] = fn
}
