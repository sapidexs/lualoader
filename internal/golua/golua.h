#ifndef golua
#define golua

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"

void* Null();

void LuaLoad(lua_State* L);

int LuaL_dofile(lua_State* L, const char* filename);

int RunEnableFunc(lua_State* L);

int RunDisableFunc(lua_State* L);

typedef struct {
    char* name;
    int ref;
} Handler;

void FreeHandlerArray(Handler* ha, int n);

Handler* GetHandler(lua_State* L, int* len);

typedef struct {
    char** data;
    int size;
} StrArr;

typedef struct {
    int* array;
    int size;
} IntArr;

typedef struct {
    char* key;
    StrArr* array;
} StrArrPair;

typedef struct {
    StrArrPair** array;
    int size;
} StrArrPairArr;

typedef struct {
    int status;
    char* body;
    StrArrPairArr* header;
} CallResult;

int SizeofStrArrPairPtr();

StrArr* NewStrArr(int len);

void StrArrSetI(StrArr* sa, int i, const char* d);

void FreeStrArr(StrArr* sa);

IntArr* NewIntArr(int len);

void IntArrSetI(IntArr* ia, int i, int v);

void FreeIntArr(IntArr* ia);

StrArrPair* NewStrArrPair(const char* k, int len);

void FreeStrArrPair(StrArrPair* sap);

StrArrPairArr* NewStrArrPairArr(int len, StrArr* kArr, IntArr* lArr);

void FreeStrArrPairArr(StrArrPairArr* sapa);

void FreeCallResult(CallResult* cr);

CallResult* CallFuncRef(lua_State* L, int ref, const char* method, const char* proto, const char* body, const char* host, const char* remoteAddr, const char* requestURI, StrArrPairArr* header, StrArrPairArr* form, StrArrPairArr* postForm, StrArrPairArr* trailer, StrArrPairArr* urlQuery);

#endif