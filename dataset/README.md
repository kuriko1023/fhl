# 数据集

整合了四个数据集

- https://github.com/chinese-poetry/chinese-poetry (f2ec391) 只并入了花间集部分
- https://github.com/javayhu/poetry (90b53bf)
- https://github.com/yxcs/poems-db (81a6fc7)
- https://github.com/snowtraces/poetry-source (7fd1315)
- https://github.com/Werneror/Poetry (0039c69)

数据集暂时还未作为子模块加入仓库，需要先将它们 clone 到目录 **10-sources** 下，并依次命名为数字 **1** 至 **5**。

数据处理步骤：
1. 将源数据集存放至 **10-sources** 下，同时将额外的内容放入 **11-extra.txt**；
2. 执行 `node merge.js`，获得以下两份文件：
    - **1a-all.txt**，包含所有格式正确、无怪异字符的条目；
    - **1b-err.txt**，包含所有不符合要求而被过滤的条目；
3. 有些诗词存在多个版本，若这些版本不全正确，将正确的那些保存在 **20-accept.txt**；
4. 将名句列表保存在 **21-curate.txt**；
5. 在 **dedup/src** 目录下执行 `cargo run --release`，获得以下两份文件：
    - **2a-dups.txt**，包含所有重复的篇目（可根据这里的信息前往更新 **20-accept.txt**）；
    - **2b-dedup.txt**，包含去重所得的最终数据集，也即服务器使用的数据集。

清理所有中间结果的命令：`rm [0-9][a-z]-*`。
