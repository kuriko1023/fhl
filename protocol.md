## 前后端交互协议

1. 根据登录方式，获取一个登录码 code
  - <small>※ 对于微信小程序，获取方式是 `wx.login()`【uni-app 封装为 `uni.login()`】</small>
  - <small>※ 对于网页版，形式为字符串 `webpage_login_` 后面拼接自定义用户名</small>
  - <small>※ 注意登录码有时是一次性的，不同于玩家 id</small>
  - <small>※ 玩家不用注册，首次登录即加入数据库</small>
  - <small style='color: #48f'>※ 测试期间可以用感叹号“!”加上任意 id 作为登录码，如用 **<span style='color: #e94'>!kuriko</span>** 作为登录码对应 id 为 **<span style='color: #e94'>kuriko</span>** 的玩家</small>
2. 与服务器建立一个 WebSocket 连接 
  - `wss://fhxl.ayu.land/channel/<room>/<code>`
  - <small>※ \<room\> 为房主的 id，另外用 **my** 表示自己的房间<span style='color: #aaa'>（因为用登录码并不能直接取得 id）</span></small>
  - <small style='color: #48f'>※ 例如　wss://fhxl.ayu.land/channel/<span style='color: #e94'>my</span>/<span style='color: #e94'>!kuriko</span><br>　　　→ 以 **<span style='color: #e94'>kuriko</span>** 身份进入自己的房间</small>
  - <small style='color: #48f'>※ 例如　wss://fhxl.ayu.land/channel/<span style='color: #e94'>kuriko</span>/<span style='color: #75f'>!ayuu</span><br>　　　→ 以 **<span style='color: #75f'>ayuu</span>** 身份进入 **<span style='color: #e94'>kuriko</span>** 的房间</small>
3. 连接建立后，所有消息均在其上以 JSON 格式传递

### 消息类型

#### 房间状态

【→ host】【→ guest】【→ spectators】
```json
{
  "type": "room_status",
  "host": "kuriko",
  "host_avatar": "kuriko/say23z",
  "host_status": "absent",
  "guest": "ayuu",
  "guest_avatar": "ayuu/say24a"
  
}
```

连接建立后的第一条消息；当有人点击坐下、进入房间或离开房间时也会发送。

`host_status` 表示房主的状态，`absent` 表示不在房间内，`present` 表示在房间内但未坐下，`ready` 表示在房间内且坐下。

`guest` 表示已经坐下的客人的名字。若没有客人坐下，则为 null。

`host_avatar` 和 `guest_avatar` 表示双方头像的标识符，这个值前接 `/avatar/` 即为头像的 URL，如本例中房主的头像为 `/avatar/kuriko/say23z`。（其中的 `say23z` 是 base 36 编码的时间戳。）

#### 坐下

【host →】【guest →】
```json
{
  "type": "ready"
}
```

服务端对此消息返回一条「房间状态」消息。

#### 进入选题

房主点击开始游戏按钮后，进入选题界面，此时即进入游戏，房间占用。

【host →】
```json
{
  "type": "start_generate"
}
```

【→ guest】
```json
{
  "type": "generate_wait"
}
```

#### 选题

房主点击选题按钮

【host →】
```json
{  
  "type": "generate",
  "mode": "A",
  "size": 0
}
```

服务端返回一道题目。

【→ host】
```json
{
  "type": "generated",
  "mode": "A",
  "size": 0,
  "subject": "草木",
  "confirm": false
}
```

换一换

【host →】
```json
{  
  "type": "generate",
  "mode": "C",
  "size": 3
}
```
【→ host】
```json
{
  "type": "generated",
  "mode": "C",
  "size": 3,
  "subject": "烟 见 横/头 春 叶 一 雨 与 依 青 绝 前 初 曾 君 秋风 无人 落日/0000000000000000",
  "confirm": false
}
```

#### 确认选题

首先，房主确认选题

【host →】
```json
{  
  "type": "confirm_subject"
}
```
【→ host】【→ guest】
```json
{
  "type": "generated",
  "mode": "C",
  "size": 3,
  "subject": "春 日 鱼/出 水 夕 夜 画 年 红 来 休 开 静 说 今 归去 草木 一片/0000000000000000",
  "confirm": false
}
```

然后等待客人阅读规则、确认题目。确认后即进入游戏。

【guest →】
```json
{  
  "type": "confirm"
}
```

#### 游戏状态

开始游戏后，服务端会向双方发送第一个「游戏状态」消息。在选题界面接收到「游戏状态」消息时，即应跳转到游戏状态页面。

【→ host】【→ guest】
```json
{  
  "type": "game_status",

  "mode": "A",  // 梦笔生花
  "subject": "花",
  "history": [
    "花谢花飞花满天/024"
  ],

  "mode": "B",  // 走马观花
  "subject": "春花秋月何时了/3",
  "history": ["..."],

  "mode": "C",  // 天女散花
  "subject": "古 梦 雁/长 舟 送 寄 事 神 不 生 西风 多少/1000010011",
  "history": ["..."],

  "mode": "D",  // 雾里看花
  "subject": "二 西/三十 故人/01/01",
  "history": [
    "西出阳关无故人/0/56"
  ],

  "host_timer": 60000,
  "guest_timer": 60000,
  "turn_timer": 60000
}
```

- `mode` 表示游戏模式，A/B/C/D
- `subject` 表示游戏题目和当前进程，含义随游戏模式不同
- `history` 为此前所有提交过的句子（游戏开始的情况下为空数组）。偶数下标（从 0 开始）对应房主，奇数下标对应客人。通过此数组的长度可以推算出当前轮到的玩家。
- `turn_timer` 表示当前轮剩余的时限；`host_timer` 和 `guest_timer` 表示双方剩余的累积加时。客户端应当根据这些值校准计时器。

#### 提交答案

【active player →】
```json
{
  "type": "answer",
  "text": "春水碧于天 画船听雨眠"
}
```

答案中的句子分隔统一使用空格，不可出现其他标点符号。

服务端对此消息返回一条「游戏状态」或「无效答案」消息。

**游戏状态消息**形式如下：

【→ host】【→ guest】
```json
{
  "type": "game_update",
  "text": "千古兴亡多少事/1/45",
  // 用 base 36 标注出高亮字下标（从 0 起，包括空格）；
  // 对于模式 C「天女散花」，固定字在前，可消除字在后；
  // 对于模式 D「雾里看花」，第一组字在前，第二组字在后
  "update": "",  // A: 没有更新
  "update": "1",   // B: 表示接下来轮到的玩家需要飞的字
  "update": "4",   // C: 用到的可消除字
  "update": "0,4", // D: 左边和右边用到的字
  "host_timer": 20,
  "guest_timer": 20
}
```

任意一方提交符合规则的答案时均会向双方发送此消息。

- `text` 表示本次提交的句子
- `update` 表示游戏题目如何更新，含义随游戏模式不同
- `host_timer` 和 `guest_timer` 表示对方剩余的累积加时。客户端应当根据这些值校准计时器。

**无效答案消息**形式如下：

【→ active player】
```json
{
  "type": "invalid_answer",
  "reason": "不审题"
}
```

只会向提交的一方发送。

`reason` 的取值有如下几种：
- `不审题`：不匹配剩余的关键词
- `复读机`：至少一句与历史中的一句重复
- `没背熟`：诗词库中存在相似但不一致的句子
- `大文豪`：诗词库中没有找到相似的句子
- `捣浆糊`：太短（&lt; 4 字，不包括空格）
- `碎碎念`：太长（&gt; 21 字，不包括空格）

有多个原因时，返回列表中最靠前的一个

#### 游戏结束

【→ host】【→ guest】
```json
{
  "type": "end_status",
  "winner": -1,
  "mode": "C",
  "subject": "烟 见 横/头 春 叶 一 雨 与 依 青 绝 前 初 曾 君 秋风 无人 落日/0001000000100000",
  "history": [
    "一江烟水照晴岚/2/0",
    "人生若只如初见/1/5",
  ]
}
```

游戏结束时发送此消息。跳转到游戏结束页面。
`winner` 为 1（自己胜）、-1（对方胜）或 0（平局，完成整道题目）。

游戏结束时，房间立刻释放，但是玩家仍然视为留在房间内（未坐下状态）。玩家浏览游戏结果后，可以回到房间界面。
