# Lucy-QOnebot 文档

Lucy | Ver.2.0 — 基于 OneBot + ZeroBot 的 QQ 机器人（HiMoYo 版本）。

- 项目仓库: [Lucy-QOnebot](https://github.com/MoYoez/Lucy-QOnebot)
- 说明书: https://lucy.lemonkoi.one
- Copyright © 2021-2024 FloatTech. All Rights Reserved.

## 文档索引

| 文档 | 说明 |
|------|------|
| [功能简介](./features.md) | 现有插件与功能一览 |

## 运行要求

- Go 1.x
- 根目录 `filter.json` 配置白名单群号
- 根目录 `.env`（如使用 godotenv）
- OneBot 兼容的 WS 服务（如 go-cqhttp）运行于 `ws://127.0.0.1:6700`

## 配置说明

- **filter.json**: 白名单群列表，启动时由 `box/whitelist` 读取。
- **SuperUsers**: 在 `main.go` 中配置超级用户 QQ 号，用于 testkey、群管、自检、whitelist 等。
