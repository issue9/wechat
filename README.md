wechat
[![Build Status](https://travis-ci.org/issue9/wechat.svg?branch=master)](https://travis-ci.org/issue9/wechat)
[![codecov](https://codecov.io/gh/issue9/wechat/branch/master/graph/badge.svg)](https://codecov.io/gh/issue9/wechat)

微信接口，测试中，勿用！

## 目录结构

```
|--- common 公众号用到的公用包
|     |
|     +------ result 表示微信的各类返回信息
|
|---- mp 公众号的相关接口
|     |
|     +----- message 消息管理
|     |
|     +----- template 模板功能
|     |
|     +----- jssdk jssdk 相关的功能
|
+---- pay 支付接口
|     |
|     +--- unifiedorder 统一支付接口
|     |
|     +--- refund 退款接口
|     |
|     +--- notify 支付通知接口
|
|---- weapp 小程序相关功能
|     |
|     +--- auth 验证
|     |
|     +--- template 模板
```
