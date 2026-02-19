# 功能简介

以下为当前主程序已启用的功能模块简介。

## 功能模块 (functions/)

| 模块 | 说明 |
|------|------|
| **choose** | 帮助做选择：触发句式如「是 A 还是 B」，Lucy 随机选一个。 |
| **daily** | 每日刷新：签到、运势等，与 score 积分、fortune 运势结合。 |
| **group** | 群组相关：群内反应、群事件处理等。 |
| **interaction** | 互动：如「叫我 XXX」保存昵称、戳一戳、随机回复等。 |
| **mai** (maidx) | 中二节奏 (MaiDX) 查分、查歌、绑定等。 |
| **manager** | 群管：禁言/解禁、全员禁言、升/取消管理、改名片/头衔、踢人、欢迎语/告别辞等。 |
| **nsfw** | NSFW 相关（默认关闭）：示例为 Hello/World 回复。 |
| **pgr** (phigros) | Phigros：绑定、查分、查曲等，依赖 PhigrosUnlimitedAPI 与本地/远程服务。 |
| **score** | 积分/签到系统：签到得柠檬片、等级、抢柠檬、转账、菠菜等。 |
| **setu** | 来图：如「来份二次元」「来份星空」等，带频率限制与自动撤回。 |
| **simai** | Simai 预渲染词典，提升对话/回复的匹配与表现。 |
| **slash** | Slash 风格指令（默认关闭）：如 `/rua`、`/xxx [CQ:at]` 等趣味互动。 |
| **tools** | 工具：超级用户可用「检查身体/自检/系统状态」查看 CPU/内存/磁盘。 |
| **whitelist** | 白名单管理：超级用户通过 `!whitelist <群号>` 动态添加白名单群。 |
| **wife** | 娶/嫁玩法：群内娶/嫁、NTR、分手等，与积分（柠檬片）联动。 |

## 公共组件 (box/)

| 路径 | 说明 |
|------|------|
| **box/whitelist** | 读取 `filter.json`，提供白名单群列表，供 main 与各插件使用。 |
| **box/notify** | 提供 Banner 等说明文案，用于 `.help` / `/help`。 |
| **box/break** | 字符串截断/分词工具，供 simai、mai、daily 等使用。 |
| **box/setname** | 用户昵称存储与读取，供 interaction、simai 等使用。 |
| **box/event** | 等待下一条消息等事件封装，供 interaction 使用。 |
| **box/coins** | 积分/签到数据库与逻辑，供 score、wife 等使用。 |
| **box/draw** | 绘图相关（字体、边框文字、图片处理等），供 daily、score 等使用。 |
| **box/emoji** | 表情/ emoji 处理，供 daily 等使用。 |

## 全局行为（main.go）

- 仅白名单内群聊会处理消息，其余群直接 `Block`。
- `.help` / `/help`：回复 Banner 说明。
- `testkey`（仅 SuperUser、OnlyToMe）：测试配置并回显当前白名单群列表。
- 所有通过 `_ "github.com/MoYoez/Lucy-QOnebot/functions/..."` 引入的插件由 ZeroBot 自动注册并响应各自触发器。

## 废弃代码

已不再在主程序中启用的旧代码已移至项目根目录下 **deprecated/** 文件夹，详见 [deprecated/README.md](../deprecated/README.md)。
