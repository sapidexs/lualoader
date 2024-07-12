#ifndef golua
#define golua

#include <stdlib.h>

#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"

const char* Lua_tostring(lua_State* L, int idx);

lua_Integer Lua_tointeger(lua_State* L, int idx);

int LuaL_dofile(lua_State* L, const char* filename);

int Lua_pcall(lua_State* L, int n, int r, int f);

#endif