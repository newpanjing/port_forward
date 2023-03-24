# go语言写的一个端口转发程序

ChatGPT编写的，我只是加个两行flag解析参数的代码

## 用法

```bash
./port_forward --help

Usage of ./port_forward:
  -listen string
        listen address (default "0.0.0.0:9999")
  -target string
        target address (default "127.0.0.1:8899")
```