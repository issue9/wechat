# wechat

微信接口，测试中，勿用！



### 目录结构

```
+---- common
|
|---- mp 公众号的相关接口
|     |
|     +--- common 公众号用到的公用包
|     |      |
|     |      +------ result 表示微信的各类返回信息
|     |      |
|     |      +------ config 微信的基本配置项
|     |      |
|     |      +------ token access_token 中控服务器的实现
|     |
|     +----- message 消息管理
|     |
|     +----- template 模板功能
|     |
|     +----- jssdk jssdk 相关的功能
|
+---- pay 支付接口
      |
      +--- unifiedorder 统一支付接口
      |
      +--- refund 退款接口
      |
      +--- notify 支付通知接口
```
