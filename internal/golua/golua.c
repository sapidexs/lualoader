#include "golua.h"
#include "lua.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void* Null() {
    return NULL;
}

void LuaLoad(lua_State* L) {
    luaL_openlibs(L);
}

int LuaL_dofile(lua_State* L, const char* filename) {
    int ret = luaL_dofile(L, filename);
    if (ret != LUA_OK) {
        printf("Lua Error: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
    }
    return ret;
}

int RunEnableFunc(lua_State* L) {
    int t = lua_getglobal(L,"Enable");
    if (t != LUA_TFUNCTION) return LUA_ERRRUN;
    int ret = lua_pcall(L, 0, 0, 0);
    if (ret != LUA_OK) {
        printf("Lua Error: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
    }
    return ret;
}

int RunDisableFunc(lua_State* L) {
    int t = lua_getglobal(L,"Disable");
    if (t != LUA_TFUNCTION) return LUA_ERRRUN;
    int ret = lua_pcall(L, 0, 0, 0);
    if (ret != LUA_OK) {
        printf("Lua Error: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
    }
    return ret;
}

Handler* NewHandlerArray(int n) {
    Handler* ha;
    ha = (Handler*)malloc(sizeof(Handler)*n);
    return ha;
}

void FreeHandlerArray(Handler* ha, int n) {
    for (int i = 0; i < n; ++i) free(ha[i].name);
    free(ha);
}

Handler* GetHandler(lua_State* L, int* len) {
    int t = lua_getglobal(L,"HandlerTable");
    if (t != LUA_TTABLE) {
        lua_pop(L, 1);
        printf("Lua-C Error: HandlerTable is not a table");
        return NULL;
    }
    int length = 0;
    lua_pushnil(L);
    while (lua_next(L, -2) != 0) {
        ++length;
        lua_pop(L, 1);
    }
    *len = length;
    Handler* handlerArray = NewHandlerArray(length);
    int i = 0;
    lua_pushnil(L);
    while (lua_next(L, -2) != 0) {
        if (lua_isstring(L, -2) == 0) {
            lua_pop(L, 3);
            FreeHandlerArray(handlerArray, length);
            printf("Lua-C Error: HandlerTable's key is not a string");
            return NULL;
        }
        const char* name = lua_tostring(L, -2);
        if (name == NULL) {
            lua_pop(L, 3);
            FreeHandlerArray(handlerArray, length);
            printf("Lua-C Error: cannot get HandlerTable's key");
            return NULL;
        }
        const int ref = luaL_ref(L, LUA_REGISTRYINDEX);
        handlerArray[i].name = (char*)malloc(strlen(name)+1);
        strcpy(handlerArray[i].name, name);
        handlerArray[i].ref = ref;
        ++i;
    }
    lua_pop(L, 1);
    return handlerArray;
}

// -----------------------------------------------------------------------------------

int SizeofStrArrPairPtr() {
    return sizeof(StrArrPair*);
}

StrArr* NewStrArr(int len) {
    StrArr* sa = (StrArr*)malloc(sizeof(StrArr));
    sa->size = len;
    sa->data = (char**)malloc(sizeof(char*)*len);
    for (int i = 0; i < len; ++i) sa->data[i] = NULL;
    return sa;
}

void StrArrSetI(StrArr* sa, int i, const char* d) {
    if (i < 0 || i >= sa->size) return;
    sa->data[i] = strdup(d);
}

void FreeStrArr(StrArr* sa) {
    for (int i = 0; i < sa->size; ++i) free(sa->data[i]);
    free(sa->data);
    free(sa);
}

IntArr* NewIntArr(int len) {
    IntArr* ia = (IntArr*)malloc(sizeof(IntArr));
    ia->size = len;
    ia->array = (int*)malloc(sizeof(int)*len);
    for (int i = 0; i < len; ++i) ia->array[i] = 0;
    return ia;
}

void IntArrSetI(IntArr* ia, int i, int v) {
    if (i < 0 || i >= ia->size) return;
    ia->array[i] = v;
}

void FreeIntArr(IntArr* ia) {
    free(ia->array);
    free(ia);
}

StrArrPair* NewStrArrPair(const char* k, int len) {
    StrArrPair* sap = (StrArrPair*)malloc(sizeof(StrArrPair));
    sap->key = strdup(k);
    sap->array = NewStrArr(len);
    return sap;
}

void FreeStrArrPair(StrArrPair* sap) {
    free(sap->key);
    FreeStrArr(sap->array);
    free(sap);
}

StrArrPairArr* NewStrArrPairArr(int len, StrArr* kArr, IntArr* lArr) {
    StrArrPairArr* sapa = (StrArrPairArr*)malloc(sizeof(StrArrPairArr));
    sapa->size = len;
    sapa->array = (StrArrPair**)malloc(sizeof(StrArrPair*)*len);
    for (int i = 0; i < len; ++i) sapa->array[i] = NewStrArrPair(kArr->data[i], lArr->array[i]);
    return sapa;
}

void FreeStrArrPairArr(StrArrPairArr* sapa) {
    for (int i = 0; i < sapa->size; ++i) FreeStrArrPair(sapa->array[i]);
    free(sapa->array);
    free(sapa);
}

CallResult* NewCallResult(int headerlen, StrArr *kArr, IntArr *lArr) {
    CallResult* cr = (CallResult*)malloc(sizeof(CallResult));
    cr->status = 0;
    cr->body = NULL;
    cr->header = NewStrArrPairArr(headerlen, kArr, lArr);
    return cr;
}

void FreeCallResult(CallResult* cr) {
    FreeStrArrPairArr(cr->header);
    free(cr->body);
    free(cr);
}

CallResult* CallFuncRef(lua_State* L, int ref, const char* method, const char* proto, const char* body, const char* host, const char* remoteAddr, const char* requestURI, StrArrPairArr* header, StrArrPairArr* form, StrArrPairArr* postForm, StrArrPairArr* trailer, StrArrPairArr* urlQuery) {
    lua_rawgeti(L, LUA_REGISTRYINDEX, ref);
    // new table
    lua_newtable(L);
    // method
    lua_pushstring(L, method);
    lua_setfield(L, -2, "method");
    // proto
    lua_pushstring(L, proto);
    lua_setfield(L, -2, "proto");
    // body
    lua_pushstring(L, body);
    lua_setfield(L, -2, "body");
    // host
    lua_pushstring(L, host);
    lua_setfield(L, -2, "host");
    // remoteAddr
    lua_pushstring(L, remoteAddr);
    lua_setfield(L, -2, "remoteAddr");
    // requestURI
    lua_pushstring(L, requestURI);
    lua_setfield(L, -2, "requestURI");
    // header
    lua_newtable(L);
    for (int i = 0; i < header->size; ++i) {
        lua_newtable(L);
        for (int j = 0; j < header->array[i]->array->size; ++j) {
            lua_pushnumber(L, j+1);
            lua_pushstring(L, header->array[i]->array->data[j]);
            lua_settable(L, -3);
        }
        lua_setfield(L, -2, header->array[i]->key);
    }
    lua_setfield(L, -2, "header");
    // form
    lua_newtable(L);
    for (int i = 0; i < form->size; ++i) {
        lua_newtable(L);
        for (int j = 0; j < form->array[i]->array->size; ++j) {
            lua_pushnumber(L, j+1);
            lua_pushstring(L, form->array[i]->array->data[j]);
            lua_settable(L, -3);
        }
        lua_setfield(L, -2, form->array[i]->key);
    }
    lua_setfield(L, -2, "form");
    // postForm
    lua_newtable(L);
    for (int i = 0; i < postForm->size; ++i) {
        lua_newtable(L);
        for (int j = 0; j < postForm->array[i]->array->size; ++j) {
            lua_pushnumber(L, j+1);
            lua_pushstring(L, postForm->array[i]->array->data[j]);
            lua_settable(L, -3);
        }
        lua_setfield(L, -2, postForm->array[i]->key);
    }
    lua_setfield(L, -2, "postForm");
    // trailer
    lua_newtable(L);
    for (int i = 0; i < trailer->size; ++i) {
        lua_newtable(L);
        for (int j = 0; j < trailer->array[i]->array->size; ++j) {
            lua_pushnumber(L, j+1);
            lua_pushstring(L, trailer->array[i]->array->data[j]);
            lua_settable(L, -3);
        }
        lua_setfield(L, -2, trailer->array[i]->key);
    }
    lua_setfield(L, -2, "trailer");
    // urlQuery
    lua_newtable(L);
    for (int i = 0; i < urlQuery->size; ++i) {
        lua_newtable(L);
        for (int j = 0; j < urlQuery->array[i]->array->size; ++j) {
            lua_pushnumber(L, j+1);
            lua_pushstring(L, urlQuery->array[i]->array->data[j]);
            lua_settable(L, -3);
        }
        lua_setfield(L, -2, urlQuery->array[i]->key);
    }
    lua_setfield(L, -2, "urlQuery");
    // get a table and return a table
    int ret = lua_pcall(L, 1, 1, 0);
    if (ret != LUA_OK) {
        printf("Lua Error: %s\n", lua_tostring(L, -1));
        lua_pop(L, 1);
        return NULL;
    }
    // get ret.header
    int t = lua_getfield(L, -1, "header");
    if (t != LUA_TTABLE) {
        lua_pop(L, 2);
        printf("Lua-C Error: header is not a table");
        return NULL;
    }
    // get length of ret.header and ret.header.kv
    int length = 0;
    lua_pushnil(L);
    while (lua_next(L, -2) != 0) {
        // check type
        if (lua_istable(L, -1) == 0) {
            lua_pop(L, 4);
            printf("Lua-C Error: header's inner value is not a table");
            return NULL;
        }
        if (lua_isstring(L, -2) == 0) {
            lua_pop(L, 4);
            printf("Lua-C Error: header's inner key is not a string");
            return NULL;
        }
        ++length;
        lua_pop(L, 1);
    }
    StrArr* kArr = NewStrArr(length);
    IntArr* lArr = NewIntArr(length);
    length = 0;
    lua_pushnil(L);
    while (lua_next(L, -2) != 0) {
        int l = 0;
        lua_pushnil(L);
        while (lua_next(L, -2) != 0) {
            ++l;
            lua_pop(L, 1);
        }
        const char* key = lua_tostring(L, -2);
        if (key == NULL) {
            lua_pop(L, 4);
            FreeStrArr(kArr);
            FreeIntArr(lArr);
            printf("Lua-C Error: cannot get header's inner key");
            return NULL;
        }
        StrArrSetI(kArr, length, key);
        IntArrSetI(lArr, length, l);
        ++length;
        lua_pop(L, 1);
    }
    // define call result
    CallResult* cr = NewCallResult(length, kArr, lArr);
    FreeStrArr(kArr);
    FreeIntArr(lArr);
    // get keys and values in ret.header
    length = 0;
    lua_pushnil(L);
    while (lua_next(L, -2) != 0) {
        int l = 0;
        lua_pushnil(L);
        while (lua_next(L, -2) != 0) {
            const char* value = lua_tostring(L, -1);
            if (value == NULL) {
                lua_pop(L, 6);
                FreeCallResult(cr);
                printf("Lua-C Error: cannot get header's inner value");
                return NULL;
            }
            StrArrSetI(cr->header->array[length]->array, l, value);
            ++l;
            lua_pop(L, 1);
        }
        ++length;
        lua_pop(L, 1);
    }
    // pop ret.header
    lua_pop(L, 1);
    // get ret.status
    t = lua_getfield(L, -1, "status");
    if (t != LUA_TNUMBER) {
        lua_pop(L, 2);
        FreeCallResult(cr);
        printf("Lua_C Error: status is not a number");
        return NULL;
    }
    cr->status = lua_tonumber(L, -1);
    // pop ret.status
    lua_pop(L, 1);
    // get ret.body
    t = lua_getfield(L, -1, "body");
    if (lua_isstring(L, -1) == 0) {
        lua_pop(L, 2);
        FreeCallResult(cr);
        printf("Lua_C Error: body is not a string");
        return NULL;
    }
    const char* retbody = lua_tostring(L, -1);
    if (retbody == NULL) {
        lua_pop(L, 2);
        FreeCallResult(cr);
        printf("Lua_C Error: cannot get body");
        return NULL;
    }
    cr->body = (char*)malloc(strlen(retbody)+1);
    strcpy(cr->body, retbody);
    // pop ret.body
    lua_pop(L, 1);
    // pop ret
    lua_pop(L, 1);
    return cr;
}