# Deprecated 废弃代码

本目录存放已不再在主程序中启用的旧功能代码，仅作归档参考。

| 路径 | 说明 |
|------|------|
| `functions/magic` | Magic 占位插件，未实现具体逻辑 |
| `functions/reborn` | 重生/重开梗（随机国家与性别），已从 main 中移除 |
| `box/ticket` | 疲劳值/Token 限流工具，当前无引用 |

如需重新启用，将对应包以 `_ "github.com/MoYoez/Lucy-QOnebot/..."` 形式在 `main.go` 中 import 即可（需同时将代码移回原目录或修改 import 路径）。
