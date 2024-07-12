# Debian12.6交叉编译lua和lualoader(CGO项目)

## lualoader是什么

+ lualoader是我写的一个简单项目，主体为go语言，通过CGO进行与原生lua的绑定。项目支持使用lua来处理简单的http请求（如表单推送，URL参数），不支持复杂http请求（如文件上传）。[项目地址](https://github.com/sapidexs/lualoader)

## 如何进行交叉编译

### lua

1. 首先下载 `gcc-riscv64-linux-gnu` 工具链。 `sudo apt install gcc-riscv64-linux-gnu` 。
2. 前往[Lua官方下载站](https://lua.org/download.html)下载最新的lua源代码，并解压。进入lua目录，执行 `make CC=riscv64-linux-gnu-gcc` 。
3. 此时 `src/` 目录下已经生成了 `lua` 和 `luac` 两个可执行文件。通过 `file` 命令查看，我们发现其需要 `/lib/ld-linux-riscv64-lp64d.so.1` 库。我们可以在 `/usr/riscv64-linux-gnu/lib/` 中找到这个库，将库复制到Duo系统的 `/lib` 下。
4. 将 `lua` 和 `luac` 两个可执行文件复制进Duo系统，测试执行。我们发现还需要 `libm.so.6` 和 `libc.so.6` 。这些同样可以在 `/usr/riscv64-linux-gnu/lib/` 下找到。我们将其复制进Duo系统的 `/lib` 下。
5. 再次执行，发现一切正常了。至此lua移植教程结束，很简单吧～

### lualoader

1. lualoader是由go语言编写的，所以我们需要先搭建好go编译环境。过程比较简单，这里不再赘述。
2. 使用 `git clone https://github.com/sapidexs/lualoader` 来下载lualoader的源码。注意**不需要**下载submodule。我们将之前编译好的lua项目 `src/` 下除了 `Makefile` ， `lua` 和 `luac` 之外的全部文件复制到 `lualoader/lua/` 文件夹下。
3. 在 `lualoader/` 文件夹下执行 `CGO_ENABLED=1 CC=riscv64-linux-gnu-gcc GOOS=linux GOARCH=riscv64 go build -ldflags "-s -w" .` 即可编译项目（怎么样，go的交叉编译很简单吧）。命令中 `-ldflags "-s -w"` 用于去除调试信息，缩小可执行文件体积，可以不写。
4. 将 `lualoader` 拷贝进Duo系统，尝试执行，能跑（因为需要的库之前已经复制进去了）但是报错 `fatal error: out of memory allocating heap arena map` （空间不足）。此时我们可以遵循[教程](https://milkv.io/zh/docs/duo/getting-started/swap)开启Duo的Swap空间。
5. 设置完Swap之后，再次执行 `lualoader` ，发现已经可以正常启动了。我们可以按照[README](https://github.com/sapidexs/lualoader/blob/main/README.md)来配置插件，并且访问网页测试看看效果。至此lualoader移植教程结束，也很简单吧～

## 作者：odorajbotoj 编辑时间：2024-07-12
