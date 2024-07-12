# lualoader

use lua to handle simple http requests

## 安装

+ Use Lip. `lip install github.com/sapidexs/lualoader` .

## 编译

### Linux amd64

1. 确保正确clone了submodule。进入lua目录并执行make。 `cd lua && make && cd ..` 。
2. 执行 `CGO_ENABLED=1 go build .` 。

### Windows amd64

1. 下载 `mingw-w64` 编译套件，并添加至环境变量。
2. 下载[dlfcn-win32](https://github.com/dlfcn-win32/dlfcn-win32)项目并编译，将 `libdl.a` 放入 `/path/to/mingw-w64/lib/` 下，将 `src/dlfcn.h` 放入 `/path/to/mingw-w64/include/` 下。
3. 从[lua官网下载页](https://www.lua.org/download.html)下载lua并编译（ `mingw32-make.exe mingw` ），将 `src/` 下除 `Makefile` `lua.exe` `luac.exe` 外所有文件放入本项目 `lua/` 文件夹下。
4. 设置 `CGO_ENABLED=1` 。
5. 执行 `go build .` 。

### Linux amd64 (Debian 12.6) 交叉编译至 RISC-V 64 (MilkV-Duo)

+ 详见 `riscv64-build.md` 。

## 运行

+ 直接运行生成的 `lualoader` 可执行文件即可。

## 配置

+ `config.json` 为配置文件，首次启动时会自动生成。
+ `port` 为服务端口，默认 `:19130` 。

## 日志

+ `errlog.txt` 文件为错误日志，会记录运行错误。

## 关于Lua

+ 使用cgo绑定原版lua(5.4)

## 插件系统

+ 插件位于 `plugins` 目录下。每个插件有单独的目录。
+ 目录下必须包含一个 `manifest.json` 配置文件。示例内容如下。

```json
{
    "manifest_version": 1,
    "name": "helloworld",
    "entry": "main.lua",
    "plugin_version": [1, 0, 0],
    "author": "odorajbotoj",
    "description": "hello, world!"
}
```

| 键名 | 含义 |
| --- | --- |
| `manifest_version` | 无需修改 |
| `name` | 插件名称 |
| `entry` | 插件入口文件 |
| `plugin_version` | 插件版本，遵循语义化版本规范 |
| `author` | 插件作者 |
| `description` | 插件描述 |

+ 插件必须包含一个 `Enable` 函数与一个 `Disable` 函数，分别在启用和禁用时调用。这两个函数无参数，无返回值。
+ 插件必须包含一个 `HandlerTable` 表，用于注册路由。这个表以字符串为键，以函数为值。函数有一个参数，有一个返回值，均为表。

+ 关于参数表的项：

| 键名 | 含义 |
| --- | --- |
| `method` | string，请求方法 |
| `proto` | string，请求协议 |
| `body` | string，请求体 |
| `host` | string，请求目标地址 |
| `remoteAddr` | string，请求源地址 |
| `requestURI` | string，请求的URI |
| `header` | string-string table嵌套string table，请求头 |
| `form` | string-string table嵌套string table，请求附带表单 |
| `postPorm` | string-string table嵌套string table，请求附带Post表单 |
| `trailer` | string-string table嵌套string table，请求额外信息 |
| `urlQuery` | string-string table嵌套string table，请求URL参数 |

+ 关于返回值表的项：

| 键名 | 含义 |
| --- | --- |
| `status` | number，响应码 |
| `body` | string，响应体 |
| `header` | string-string table嵌套string table，响应头 |

+ [插件示例](https://github.com/sapidexs/lualoader-plugin-demo)
