# changes

- 2022-07-20 支持使用 [sonic](https://github.com/bytedance/sonic) JSON 编解码器，压测性能略有提升，没有显著提升。

压测环境 | 配置   | 接口    | easyjson | sonic
---------|--------|---------|----------|-------
本机     | 16G12C | /status | 93529    | 101120
Linux    | 64G32C | /status | 233579   | 239860
