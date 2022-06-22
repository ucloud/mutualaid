# 开发指南

本项目为微信小程序，开发前需要注册申请微信小程序账号、安装微信小程序开发者工具等前置工作，具体可参考小程序的官方文档：https://developers.weixin.qq.com/miniprogram/dev/framework/quickstart/#%E5%B0%8F%E7%A8%8B%E5%BA%8F%E7%AE%80%E4%BB%8B

## 关于依赖

项目所依赖的第三方库，已通过微信小程序开发者工具的构建功能完成构建。若在二次开发中，需要添加 npm 依赖，则必须重新执行构建。具体可参考：https://developers.weixin.qq.com/miniprogram/dev/devtools/npm.html

## 目录结构

```
.
├── app.js // 小程序主入口
├── app.json // 小程序基本配置
├── app.wxss // 小程序全局样式
├── components // 自定义组件
├── images // 图片
├── lib // 库
├── miniprogram_npm // 第三方依赖
├── package.json
├── pages // 页面
├── project.config.json
├── sitemap.json
├── utils // 工具函数
```

## 开发环境

微信开发者工具 Stable 1.05.2204250

## 自定义配置

### 地理配置说明

地理位置用于就近互助，不授权则使用上海中心位置代替。  
默认地理经纬度上海中心设置为 const SH_latitude = 31.2381;const SH_longitude = 121.4692;  
如果需要默认其它位置，在 app.js 文件中修改地址相对应的经纬度就可以。

### 文案相关配置说明

在关于页面中的相关文案信息文字位置在项目 my/my.wxss 文件中，直接修改相关文案即可。

### 小程序消息订阅提醒

本小程序提供消息订阅提醒功能，需要将 app.js 文件中，`wxSubscribeGotHelp` 和 `wxSubscribeAcceptHelp` 中的参数 `tmplIds` 改为自己账号下申请的模板 ID，如下图：

![模板ID](https://www-s.ucloud.cn/2022/06/e2e914a62229d3fdc60bb46ec9388ff1_1655866083491.jpg)
