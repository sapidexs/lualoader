package golua

/*
#cgo CFLAGS: -I${SRCDIR}/../../lua
#cgo LDFLAGS: -L${SRCDIR}/../../lua -llua -lm -ldl

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
	LuaStatePool map[string]*C.lua_State
	Mux          *http.ServeMux
)

func InitLuaStatePool() {
	LuaStatePool = make(map[string]*C.lua_State)
}

func LuaNewState() *C.lua_State {
	L := C.luaL_newstate()
	C.luaL_openlibs(L)
	return L
}

func LuaDoFile(L *C.lua_State, path string) error {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	ret := C.LuaL_dofile(L, p)
	if ret != C.LUA_OK {
		err := errors.New("Lua Error: " + C.GoString(C.Lua_tostring(L, -1)))
		C.lua_settop(L, -2)
		return err
	}
	return nil
}

func LuaCloseState(L *C.lua_State) {
	C.lua_close(L)
}

func LuaPluginRunEnableFunc(L *C.lua_State) error {
	fn := C.CString("Enable")
	defer C.free(unsafe.Pointer(fn))
	t := C.lua_getglobal(L, fn)
	if t != C.LUA_TFUNCTION {
		C.lua_settop(L, -2)
		return errors.New("Lua Error: \"Enable\" is not a function")
	}
	ret := C.Lua_pcall(L, 0, 0, 0)
	if ret != C.LUA_OK {
		err := errors.New("Lua Error: " + C.GoString(C.Lua_tostring(L, -1)))
		C.lua_settop(L, -2)
		return err
	}
	return nil
}

func LuaPluginRunDisableFunc(L *C.lua_State) error {
	fn := C.CString("Disable")
	defer C.free(unsafe.Pointer(fn))
	t := C.lua_getglobal(L, fn)
	if t != C.LUA_TFUNCTION {
		C.lua_settop(L, -2)
		return errors.New("Lua Error: \"Disable\" is not a function")
	}
	ret := C.Lua_pcall(L, 0, 0, 0)
	if ret != C.LUA_OK {
		err := errors.New("Lua Error: " + C.GoString(C.Lua_tostring(L, -1)))
		C.lua_settop(L, -2)
		return err
	}
	return nil
}

func LuaPluginGetHandler(L *C.lua_State) (map[string]C.int, error) {
	tname := C.CString("HandlerTable")
	defer C.free(unsafe.Pointer(tname))
	t := C.lua_getglobal(L, tname)
	if t != C.LUA_TTABLE {
		C.lua_settop(L, -2)
		return nil, errors.New("Lua Error: \"HandlerTable\" is not a table")
	}
	handlerTableMap := make(map[string]C.int)
	C.lua_pushnil(L)
	for C.lua_next(L, -2) != 0 {
		if C.lua_isstring(L, -2) == 0 {
			C.lua_settop(L, -4)
			return nil, errors.New("Lua Error: HandlerTable's key is not a string")
		}
		name := C.Lua_tostring(L, -2)
		if unsafe.Pointer(name) == C.NULL {
			C.lua_settop(L, -4)
			return nil, errors.New("Lua Error: cannot get HandlerTable's key")
		}
		handlerTableMap[C.GoString(name)] = C.luaL_ref(L, C.LUA_REGISTRYINDEX)
	}
	C.lua_settop(L, -2)
	return handlerTableMap, nil
}

func luaSetStrField(L *C.lua_State, k string, v string) {
	key := C.CString(k)
	defer C.free(unsafe.Pointer(key))
	value := C.CString(v)
	defer C.free(unsafe.Pointer(value))
	C.lua_pushstring(L, value)
	C.lua_setfield(L, -2, key)
}

func luaSetNumField(L *C.lua_State, k int, v string) {
	value := C.CString(v)
	defer C.free(unsafe.Pointer(value))
	C.lua_pushnumber(L, C.double(k))
	C.lua_pushstring(L, value)
	C.lua_settable(L, -3)
}

func luaSetMapStrStrArr(L *C.lua_State, keystr string, m map[string][]string) {
	C.lua_createtable(L, 0, 0)
	for k, v := range m {
		C.lua_createtable(L, 0, 0)
		for k2, v2 := range v {
			luaSetNumField(L, k2+1, v2)
		}
		key := C.CString(k)
		C.lua_setfield(L, -2, key)
		C.free(unsafe.Pointer(key))
	}
	KEY := C.CString(keystr)
	defer C.free(unsafe.Pointer(KEY))
	C.lua_setfield(L, -2, KEY)
}

func LuaPluginSetHandler(L *C.lua_State) error {
	m, err := LuaPluginGetHandler(L)
	if err != nil {
		return err
	}
	for k, v := range m {
		Mux.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			// get callback function
			C.lua_rawgeti(L, C.LUA_REGISTRYINDEX, C.longlong(v))
			// set args
			// create table
			C.lua_createtable(L, 0, 0)
			// method
			luaSetStrField(L, "method", r.Method)
			// proto
			luaSetStrField(L, "proto", r.Proto)
			// body
			bodybytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			luaSetStrField(L, "body", string(bodybytes))
			// host
			luaSetStrField(L, "host", r.Host)
			// remoteAddr
			luaSetStrField(L, "remoteAddr", r.RemoteAddr)
			// requestURI
			luaSetStrField(L, "requestURI", r.RequestURI)
			// header
			luaSetMapStrStrArr(L, "header", r.Header)
			// form
			luaSetMapStrStrArr(L, "form", r.Form)
			// postForm
			luaSetMapStrStrArr(L, "postForm", r.PostForm)
			// trailer
			luaSetMapStrStrArr(L, "trailer", r.Trailer)
			// urlQuery
			luaSetMapStrStrArr(L, "urlQuery", r.URL.Query())
			// call the function, the function gets a table and returns a table
			ok := C.Lua_pcall(L, 1, 1, 0)
			if ok != C.LUA_OK {
				w.Write([]byte("Lua Error: " + C.GoString(C.Lua_tostring(L, -1))))
				C.lua_settop(L, -2)
				return
			}
			// get ret.header
			HEADER := C.CString("header")
			defer C.free(unsafe.Pointer(HEADER))
			t := C.lua_getfield(L, -1, HEADER)
			if t != C.LUA_TTABLE {
				C.lua_settop(L, -3)
				w.Write([]byte("Lua Error: header is not a table"))
				return
			}
			C.lua_pushnil(L)
			for C.lua_next(L, -2) != 0 {
				// check value
				if C.lua_type(L, -1) != C.LUA_TTABLE {
					C.lua_settop(L, -5)
					w.Write([]byte("Lua Error: header's inner value is not a table"))
					return
				}
				// check key
				if C.lua_isstring(L, -2) == 0 {
					C.lua_settop(L, -5)
					w.Write([]byte("Lua Error: header's inner key is not a string"))
					return
				}
				// get key
				key := C.Lua_tostring(L, -2)
				if unsafe.Pointer(key) == C.NULL {
					C.lua_settop(L, -5)
					w.Write([]byte("Lua Error: cannot get header's inner key"))
					return
				}
				// get inner table
				C.lua_pushnil(L)
				for C.lua_next(L, -2) != 0 {
					// get value
					val := C.Lua_tostring(L, -1)
					if unsafe.Pointer(val) == C.NULL {
						C.lua_settop(L, -7)
						w.Write([]byte("Lua Error: cannot get header's inner value"))
						return
					}
					w.Header().Add(C.GoString(key), C.GoString(val))
					C.lua_settop(L, -2)
				}
				C.lua_settop(L, -2)
			}
			// pop ret.header
			C.lua_settop(L, -2)
			// get ret.status
			STATUS := C.CString("status")
			defer C.free(unsafe.Pointer(STATUS))
			t = C.lua_getfield(L, -1, STATUS)
			if t != C.LUA_TNUMBER {
				C.lua_settop(L, -3)
				w.Write([]byte("Lua Error: status is not a number"))
				return
			}
			statcode := C.Lua_tointeger(L, -1)
			if statcode == 0 {
				C.lua_settop(L, -3)
				w.Write([]byte("Lua Error: bad status code"))
				return
			}
			w.WriteHeader(int(statcode))
			// pop ret.status
			C.lua_settop(L, -2)
			// get ret.body
			BODY := C.CString("body")
			defer C.free(unsafe.Pointer(BODY))
			t = C.lua_getfield(L, -1, BODY)
			if C.lua_isstring(L, -1) == 0 {
				C.lua_settop(L, -3)
				w.Write([]byte("Lua Error: body is not a string"))
				return
			}
			retbody := C.Lua_tostring(L, -1)
			if unsafe.Pointer(retbody) == C.NULL {
				C.lua_settop(L, -3)
				w.Write([]byte("Lua Error: cannot get body"))
				return
			}
			w.Write([]byte(C.GoString(retbody)))
			// pop ret.body
			C.lua_settop(L, -2)
			// pop ret
			C.lua_settop(L, -2)
			// return
			return
		})
	}
	return nil
}
