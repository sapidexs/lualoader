#include "golua.h"

const char *Lua_tostring(lua_State *L, int idx) { return lua_tostring(L, idx); }

lua_Integer Lua_tointeger(lua_State *L, int idx) {
  return lua_tointeger(L, idx);
}

int LuaL_dofile(lua_State *L, const char *filename) {
  return luaL_dofile(L, filename);
}

int Lua_pcall(lua_State *L, int n, int r, int f) {
  return lua_pcall(L, n, r, f);
}
