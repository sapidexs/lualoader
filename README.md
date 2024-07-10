# lualoader

use lua to handle simple http requests

## 编译

+ TODO

## 运行

+ 直接运行生成的 `lualoader` 可执行文件即可。

## 配置

+ `config.json` 为配置文件，首次启动时会自动生成。
+ `port` 为服务端口，默认 `:19130` 。

## 日志

+ `errlog.txt` 文件为错误日志，会记录**golang侧**运行错误。

## 关于Lua

+ 使用cgo绑定原版lua

## 插件系统

+ 插件位于 `plugins` 目录下。每个插件有单独的目录。
+ 目录下必须包含一个 `manifest.json` 配置文件。示例内容如下。

```json
{
    "manifest_version": 1, // 无需修改
    "name": "helloworld", // 插件名称
    "entry": "main.lua", // 插件入口文件
    "plugin_version": [1, 0, 0], // 插件版本，遵循语义化版本规范
    "author": "odorajbotoj", // 插件作者
    "description": "hello, world!" // 插件描述
}
```

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
