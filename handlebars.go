package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"unsafe"

	"github.com/aymerick/raymond"
)

func main() {}

func GetArgumentAsString(args *C.HalonHSLArguments, pos uint64, required bool) (string, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return "", fmt.Errorf("missing argument at position %d", pos)
		} else {
			return "", nil
		}
	}
	var y *C.char
	if C.HalonMTA_hsl_value_get(x, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&y), nil) {
		return C.GoString(y), nil
	} else {
		return "", fmt.Errorf("invalid argument at position %d", pos)
	}
}

func GetArgumentAsJSON(args *C.HalonHSLArguments, pos uint64, required bool) (string, error) {
	var x = C.HalonMTA_hsl_argument_get(args, C.ulong(pos))
	if x == nil {
		if required {
			return "", fmt.Errorf("missing argument at position %d", pos)
		} else {
			return "", nil
		}
	}
	var y *C.char
	z := C.HalonMTA_hsl_value_to_json(x, &y, nil)
	defer C.free(unsafe.Pointer(y))
	if z {
		return C.GoString(y), nil
	} else {
		return "", fmt.Errorf("invalid argument at position %d", pos)
	}
}

func SetReturnValueToAny(ret *C.HalonHSLValue, val interface{}) error {
	x, err := json.Marshal(val)
	if err != nil {
		return err
	}
	y := C.CString(string(x))
	defer C.free(unsafe.Pointer(y))
	var z *C.char
	if !(C.HalonMTA_hsl_value_from_json(ret, y, &z, nil)) {
		if z != nil {
			err = errors.New(C.GoString(z))
			C.free(unsafe.Pointer(z))
		} else {
			err = errors.New("failed to parse return value")
		}
		return err
	}
	return nil
}

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

//export Halon_init
func Halon_init(hic *C.HalonInitContext) C.bool {
	return true
}

//export handlebars
func handlebars(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	template, err := GetArgumentAsString(args, 0, true)
	if err != nil {
		value := map[string]interface{}{"error": err.Error()}
		SetReturnValueToAny(ret, value)
		return
	}
	context, err := GetArgumentAsJSON(args, 1, true)
	if err != nil {
		value := map[string]interface{}{"error": err.Error()}
		SetReturnValueToAny(ret, value)
		return
	}

	var ctx map[string]interface{}
	json.Unmarshal([]byte(context), &ctx)

	result, err := raymond.Render(template, ctx)
	if err != nil {
		value := map[string]interface{}{"error": err.Error()}
		SetReturnValueToAny(ret, value)
		return
	}

	value := map[string]interface{}{"result": result}
	SetReturnValueToAny(ret, value)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	handlebars_cs := C.CString("handlebars")
	C.HalonMTA_hsl_register_function(hhrc, handlebars_cs, nil)
	C.HalonMTA_hsl_module_register_function(hhrc, handlebars_cs, nil)
	return true
}
