package code_engine

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/dop251/goja_nodejs/require"
)

var ProjectPath = "./"

func init() {
	// abs, _ := filepath.Abs(".")
	// ProjectPath = RealPath(abs)
	ProjectPath = ProgramPath()
	if strings.Contains(ProjectPath, TmpPath("")) {
		ProjectPath = RootPath()
	}
}

// RootPath Project Launch Path
func RootPath() string {
	path, _ := filepath.Abs(".")
	return RealPath(path)
}

func TmpPath(pattern ...string) string {
	p := ""
	if len(pattern) > 0 {
		p = pattern[0]
	}
	path, _ := ioutil.TempDir("", p)
	if p == "" {
		path, _ = filepath.Split(path)
	}

	if p, err := filepath.EvalSymlinks(path); err == nil {
		path = p
	}
	return RealPath(path)
}

func (vm *JsVM) sourceLoader(_ string) func(string) ([]byte, error) {
	return func(filename string) ([]byte, error) {
		// fmt.Println("sourceLoader -> filename", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			for _, v := range []string{filename, filename + ".ts"} {
				data, err = os.ReadFile(v)
				if err == nil {
					// fmt.Println("v", v)
					break
				}
			}
		}

		if err != nil {
			return nil, require.ModuleFileDoesNotExistError
		}
		ext := filepath.Ext(filename)
		if ext == ".ts" || ext == "" {
			data, err = vm.Transpile(data)
		}

		return data, err
	}
}

func pathAddSlash(path string, addSlash ...bool) string {
	if len(addSlash) > 0 && addSlash[0] && !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

func Base64EncodeString(value string) string {
	data := String2Bytes(value)
	return base64.StdEncoding.EncodeToString(data)
}

// Bytes2String bytes to string
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes string to bytes
// remark: read only, the structure of runtime changes will be affected, the role of unsafe.Pointer will be changed, and it will also be affected
func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// ReadFile ReadFile
func ReadFile(path string) ([]byte, error) {
	path = RealPath(path)
	return ioutil.ReadFile(path)
}

// ProgramPath program directory path
func ProgramPath(addSlash ...bool) (path string) {
	ePath, err := os.Executable()
	if err != nil {
		ePath = ProjectPath
	} else {
		ePath = filepath.Dir(ePath)
	}
	realPath, err := filepath.EvalSymlinks(ePath)
	if err == nil {
		ePath = realPath
	}
	path = RealPath(ePath, addSlash...)

	return
}

// RealPath get an absolute path
func RealPath(path string, addSlash ...bool) (realPath string) {
	if len(path) > 2 && path[1] == ':' {
		realPath = path
	} else {
		if len(path) == 0 || (path[0] != '/' && !filepath.IsAbs(path)) {
			path = ProjectPath + "/" + path
		}
		realPath, _ = filepath.Abs(path)
	}

	realPath = strings.Replace(realPath, "\\", "/", -1)
	realPath = pathAddSlash(realPath, addSlash...)

	return
}

func Base64DecodeString(data string) (value string, err error) {
	var dst []byte
	dst, err = base64.StdEncoding.DecodeString(data)
	if err == nil {
		value = Bytes2String(dst)
	}
	return
}

type appString interface {
	String() string
}

func ToString(i interface{}) string {
	if i == nil {
		return ""
	}
	switch value := i.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.Itoa(int(value))
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return Bytes2String(value)
	default:
		if f, ok := value.(appString); ok {
			return f.String()
		}
		return toJsonString(value)
	}
}

func toJsonString(value interface{}) string {
	jsonContent, _ := json.Marshal(value)
	jsonContent = bytes.Trim(jsonContent, `"`)
	return Bytes2String(jsonContent)
}
