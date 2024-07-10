package golua

/*
#cgo CFLAGS: -I${SRCDIR}/../../lua
#cgo !windows LDFLAGS: -L${SRCDIR}/../../lua -llua -lm -ldl
#cgo windows LDFLAGS: -L${SRCDIR}/../../lua liblua.a libm.a libdl.a

#include "golua.h"
*/
import "C"

import (
	"errors"
	"io"
	"net/http"
	"unsafe"
)

var (
	LuaStatePool map[string]*C.struct_lua_State
	Mux          *http.ServeMux
)

func InitLuaStatePool() {
	LuaStatePool = make(map[string]*C.struct_lua_State)
}

func LuaNewState() *C.struct_lua_State {
	L := C.luaL_newstate()
	C.LuaLoad(L)
	return L
}

func LuaDoFile(L *C.struct_lua_State, path string) int {
	p := C.CString(path)
	r := int(C.LuaL_dofile(L, p))
	C.free(unsafe.Pointer(p))
	return r
}

func LuaCloseState(L *C.struct_lua_State) {
	C.lua_close(L)
}

func LuaPluginRunEnableFunc(L *C.struct_lua_State) int {
	return int(C.RunEnableFunc(L))
}

func LuaPluginRunDisableFunc(L *C.struct_lua_State) int {
	return int(C.RunDisableFunc(L))
}

func LuaPluginGetHandler(L *C.struct_lua_State) error {
	var length C.int = 0
	handlerArray := C.GetHandler(L, &length)
	if unsafe.Pointer(handlerArray) == C.Null() {
		return errors.New("Error in getting handler.")
	}
	defer C.FreeHandlerArray(handlerArray, length)
	for i := 0; i < int(length); i++ {
		handler := *(*C.Handler)(unsafe.Pointer(uintptr(unsafe.Pointer(handlerArray)) + uintptr(C.sizeof_Handler*C.int(i))))
		Mux.HandleFunc(C.GoString(handler.name), func(w http.ResponseWriter, r *http.Request) {
			// args
			bodybytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			method := C.CString(r.Method)
			defer C.free(unsafe.Pointer(method))
			proto := C.CString(r.Proto)
			defer C.free(unsafe.Pointer(proto))
			body := C.CString(string(bodybytes))
			defer C.free(unsafe.Pointer(body))
			host := C.CString(r.Host)
			defer C.free(unsafe.Pointer(host))
			remoteAddr := C.CString(r.RemoteAddr)
			defer C.free(unsafe.Pointer(remoteAddr))
			requestURI := C.CString(r.RequestURI)
			defer C.free(unsafe.Pointer(requestURI))
			// Header
			headerKArr := C.NewStrArr(C.int(len(r.Header)))
			defer C.FreeStrArr(headerKArr)
			headerLArr := C.NewIntArr(C.int(len(r.Header)))
			defer C.FreeIntArr(headerLArr)
			index := 0
			for k, v := range r.Header {
				key := C.CString(k)
				C.StrArrSetI(headerKArr, C.int(index), key)
				C.IntArrSetI(headerLArr, C.int(index), C.int(len(v)))
				C.free(unsafe.Pointer(key))
				index++
			}
			header := C.NewStrArrPairArr(C.int(len(r.Header)), headerKArr, headerLArr)
			defer C.FreeStrArrPairArr(header)
			index = 0
			for _, v := range r.Header {
				index2 := 0
				for _, v2 := range v {
					value2 := C.CString(v2)
					C.StrArrSetI((*(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(header.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(index))))).array, C.int(index2), value2)
					C.free(unsafe.Pointer(value2))
					index2++
				}
				index++
			}
			// Form
			formKArr := C.NewStrArr(C.int(len(r.Form)))
			defer C.FreeStrArr(formKArr)
			formLArr := C.NewIntArr(C.int(len(r.Form)))
			defer C.FreeIntArr(formLArr)
			index = 0
			for k, v := range r.Form {
				key := C.CString(k)
				C.StrArrSetI(formKArr, C.int(index), key)
				C.IntArrSetI(formLArr, C.int(index), C.int(len(v)))
				C.free(unsafe.Pointer(key))
				index++
			}
			form := C.NewStrArrPairArr(C.int(len(r.Form)), formKArr, formLArr)
			defer C.FreeStrArrPairArr(form)
			index = 0
			for _, v := range r.Form {
				index2 := 0
				for _, v2 := range v {
					value2 := C.CString(v2)
					C.StrArrSetI((*(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(form.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(index))))).array, C.int(index2), value2)
					C.free(unsafe.Pointer(value2))
					index2++
				}
				index++
			}
			// PostForm
			postFormKArr := C.NewStrArr(C.int(len(r.PostForm)))
			defer C.FreeStrArr(postFormKArr)
			postFormLArr := C.NewIntArr(C.int(len(r.PostForm)))
			defer C.FreeIntArr(postFormLArr)
			index = 0
			for k, v := range r.PostForm {
				key := C.CString(k)
				C.StrArrSetI(postFormKArr, C.int(index), key)
				C.IntArrSetI(postFormLArr, C.int(index), C.int(len(v)))
				C.free(unsafe.Pointer(key))
				index++
			}
			postForm := C.NewStrArrPairArr(C.int(len(r.PostForm)), postFormKArr, postFormLArr)
			defer C.FreeStrArrPairArr(postForm)
			index = 0
			for _, v := range r.PostForm {
				index2 := 0
				for _, v2 := range v {
					value2 := C.CString(v2)
					C.StrArrSetI((*(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(postForm.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(index))))).array, C.int(index2), value2)
					C.free(unsafe.Pointer(value2))
					index2++
				}
				index++
			}
			// Trailer
			trailerKArr := C.NewStrArr(C.int(len(r.Trailer)))
			defer C.FreeStrArr(trailerKArr)
			trailerLArr := C.NewIntArr(C.int(len(r.Trailer)))
			defer C.FreeIntArr(trailerLArr)
			index = 0
			for k, v := range r.Trailer {
				key := C.CString(k)
				C.StrArrSetI(trailerKArr, C.int(index), key)
				C.IntArrSetI(trailerLArr, C.int(index), C.int(len(v)))
				C.free(unsafe.Pointer(key))
				index++
			}
			trailer := C.NewStrArrPairArr(C.int(len(r.Trailer)), trailerKArr, trailerLArr)
			defer C.FreeStrArrPairArr(trailer)
			index = 0
			for _, v := range r.Trailer {
				index2 := 0
				for _, v2 := range v {
					value2 := C.CString(v2)
					C.StrArrSetI((*(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(trailer.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(index))))).array, C.int(index2), value2)
					C.free(unsafe.Pointer(value2))
					index2++
				}
				index++
			}
			// URLQuery
			urlq := r.URL.Query()
			urlQueryKArr := C.NewStrArr(C.int(len(urlq)))
			defer C.FreeStrArr(urlQueryKArr)
			urlQueryLArr := C.NewIntArr(C.int(len(urlq)))
			defer C.FreeIntArr(urlQueryLArr)
			index = 0
			for k, v := range urlq {
				key := C.CString(k)
				C.StrArrSetI(urlQueryKArr, C.int(index), key)
				C.IntArrSetI(urlQueryLArr, C.int(index), C.int(len(v)))
				C.free(unsafe.Pointer(key))
				index++
			}
			urlQuery := C.NewStrArrPairArr(C.int(len(urlq)), urlQueryKArr, urlQueryLArr)
			defer C.FreeStrArrPairArr(urlQuery)
			index = 0
			for _, v := range urlq {
				index2 := 0
				for _, v2 := range v {
					value2 := C.CString(v2)
					C.StrArrSetI((*(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(urlQuery.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(index))))).array, C.int(index2), value2)
					C.free(unsafe.Pointer(value2))
					index2++
				}
				index++
			}
			// call
			callrst := C.CallFuncRef(L, handler.ref, method, proto, body, host, remoteAddr, requestURI, header, form, postForm, trailer, urlQuery)
			// return
			if unsafe.Pointer(callrst) == C.Null() {
				w.Write([]byte("Error in calling lua function"))
				return
			}
			defer C.FreeCallResult(callrst)
			for i := 0; i < int(callrst.header.size); i++ {
				hd := *(**C.StrArrPair)(unsafe.Pointer(uintptr(unsafe.Pointer(callrst.header.array)) + uintptr(C.SizeofStrArrPairPtr()*C.int(i))))
				for j := 0; j < int(hd.array.size); j++ {
					val := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(hd.array.data)) + uintptr(j)*unsafe.Sizeof(hd.array.data)))
					w.Header().Add(C.GoString(hd.key), C.GoString(val))
				}
			}
			w.WriteHeader(int(callrst.status))
			w.Write([]byte(C.GoString(callrst.body)))
		})
	}

	return nil
}
