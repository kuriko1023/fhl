# 飞花小令

在古代酒令游戏“飞花令”基础上设计的双人诗词游戏，共有四种多样的玩法。借助完善的诗词数据库，支持随机出题、检查答案。

**在线游玩** → **[fhxl.ayu.land](https://fhxl.ayu.land)**

## 玩法

“飞花令”是一种中国古代的酒令游戏，参与者轮流吟出包含特定字且格律固定的诗句。受近年《中国诗词大会》节目的启发，游戏中创造了多种不同的玩法——
- **梦笔生花**：题目为单字或二字词，玩家轮流说出带有该字或词的诗句。
- **走马观花**：题目为一句诗句（可选择 5—9 的字数），其中的字按顺序依次作为关键字，玩家轮流各自说出包含当前关键字的诗句。
- **天女散花**：题目为一组固定字词与一组可消去字词（可选择“1 词 + 10 词”或“3 词 + 16 词”），玩家轮流从两组字词中各选择一个，说出同时含有两者的诗句。每个消去词只能被选择一次。
- **雾里看花**：题目为两组字词（可选择每组 5—10 词）。玩家轮流从两组字词中各选择一个，说出同时含有两者的诗句。所有词都只能被选择一次。

我们从网络上收集了数个诗词库，合并整理并进行了一些校对工作，形成了数十万首作品的文库；从中整理出常用词句与知名诗句用于自动随机命题，细心平衡难度，以求游戏富于趣味。

## 仓库组织与细节

- **dataset/**：文库
  - 由数个诗词数据库合并、去重获得（文件 **2b-dedup.txt**）。
  - 数据处理程序由 JavaScript 和 Rust 语言实现。
  - 目前所用的此文件 SHA-2-256 摘要值为 `3bbd22542f9288f9d20c3e38a490bb404a79a9b5bbce0c0b98c6f574a8b25ec0`。
- **backend/**：服务端
  - 服务端由 Go 语言实现。
  - 服务端与客户端之间主要通过 WebSocket 通信（另有少量 HTTP），详见[交互协议文档](protocol.md)。
  - 服务端初次运行时会对文库数据作预处理（文件 **2c-precal.bin**）。
  - 关于自动出题、相似句查找等算法详情，参见源码文件 [**logic.go**](backend/logic.go)。
- **frontend/**：客户端
  - 客户端在 uni-app 框架下实现，同时适配微信小程序（已暂停运营）与网页端。
  - 使用的程序包较旧，近期的 Node.js 可能需要稍作调整才能成功构建。（具体而言，针对此文档撰写时的 LTS 版本 Node v20，需设置环境变量 `NODE_OPTIONS=--openssl-legacy-provider`。⚠️注：由于本项目中 Node 只用于构建步骤，在此选择这一做法；在其他项目中采取类似操作时，仍应谨慎考虑安全问题。）

## 许可证

源码在木兰公共许可证下分发，许可证文本见 [**COPYING.MulanPubL.md**](COPYING.MulanPubL.md)。其余资源在 [CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/) 许可证下分发。

## International Introduction

**Fei Hua Xiao Ling (飞花小令)** is a game based on an ancient Chinese drinking game of poetry, with added innovations.

Fei Hua Ling (飞花令) is a game where players recite or compose poetry given a theme, which is usually a Chinese character (it was originally fixed to be “花”, “flower”; but variations existed). The verses should also adhere to metrical rules. Inspired by the Chinese Poetry Congress television programme, a set of new gameplay rules were implemented in the game, including themes of two-syllable words, characters in verses, and pools of words.

We merged several Chinese poetry collections on the Internet and added emendations, obtaining a library of hundreds of thousands of works. We curated a set of frequent words and phrases as well as well-known verses for automatic theme generation and carefully tuned the difficulties, trying our best to deliver interesting and enjoyable games.
