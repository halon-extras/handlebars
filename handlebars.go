package main

// #cgo CFLAGS: -I/opt/halon/include
// #cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-all
// #include <HalonMTA.h>
// #include <stdlib.h>
import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/aymerick/raymond"
)

func main() {}

//export Halon_version
func Halon_version() C.int {
	return C.HALONMTA_PLUGIN_VERSION
}

//export Halon_init
func Halon_init(hic *C.HalonInitContext) C.bool {
	return true
}

func set_ret_value(ret *C.HalonHSLValue, key string, value string) {
	var ret_key *C.HalonHSLValue
	var ret_value *C.HalonHSLValue
	C.HalonMTA_hsl_value_array_add(ret, &ret_key, &ret_value)
	key_cs := C.CString(key)
	key_cs_up := unsafe.Pointer(key_cs)
	defer C.free(key_cs_up)
	value_cs := C.CString(value)
	value_cs_up := unsafe.Pointer(value_cs)
	defer C.free(value_cs_up)

	C.HalonMTA_hsl_value_set(ret_key, C.HALONMTA_HSL_TYPE_STRING, key_cs_up, 0)
	C.HalonMTA_hsl_value_set(ret_value, C.HALONMTA_HSL_TYPE_STRING, value_cs_up, 0)
}

//export handlebars
func handlebars(hhc *C.HalonHSLContext, args *C.HalonHSLArguments, ret *C.HalonHSLValue) {
	var template string
	var handlebars_cs *C.char

	var args_0 = C.HalonMTA_hsl_argument_get(args, 0)
	if args_0 != nil {
		if !C.HalonMTA_hsl_value_get(args_0, C.HALONMTA_HSL_TYPE_STRING, unsafe.Pointer(&handlebars_cs), nil) {
			set_ret_value(ret, "error", "Invalid type of \"template\" argument")
			return
		}
		template = C.GoString(handlebars_cs)
	} else {
		set_ret_value(ret, "error", "Missing required \"id\" argument")
		return
	}

	var context_str string
	var args_1 = C.HalonMTA_hsl_argument_get(args, 1)
	if args_1 != nil {
		var context_cs *C.char
		defer C.free(unsafe.Pointer(context_cs))
		var context_size_t C.size_t
		var success = C.HalonMTA_hsl_value_to_json(args_1, &context_cs, &context_size_t)
		context_str = C.GoString(context_cs)
		if !success {
			set_ret_value(ret, "error", context_str)
			return
		}
	} else {
		set_ret_value(ret, "error", "Missing required \"context\" argument")
		return
	}

	var context map[string]interface{}
	json.Unmarshal([]byte(context_str), &context)

	result, err := raymond.Render(template, context)
	if err != nil {
		set_ret_value(ret, "error", err.Error())
	}

	set_ret_value(ret, "result", result)
}

//export Halon_hsl_register
func Halon_hsl_register(hhrc *C.HalonHSLRegisterContext) C.bool {
	handlebars_cs := C.CString("handlebars")
	C.HalonMTA_hsl_register_function(hhrc, handlebars_cs, nil)
	C.HalonMTA_hsl_module_register_function(hhrc, handlebars_cs, nil)
	return true
}
