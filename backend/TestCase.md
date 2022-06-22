# 简介

本文档描述主要接口的测试用例。执行这些测试用例需要 grpcurl 和 curl。
请安装好这两个工具。

## 服务存活测试

    grpcurl -plaintext localhost:29000 list

## 下单 -- 单个商品使用优惠券

    {
      "Action": "CreateOrder",
      "account_id": 60160,
      "organization_id": 1,
      "top_organization_id": 1,
      "ZoneId": 666003001,
      "channel": 1,
      "CouponCode": "PBIXTZBX",
      "PayMethod": 1,
      "Recipient": {
          "Name": "王永庆",
          "Email": "wolf.sr@metaverse.com",
          "Mobile": "13499192345",
          "Province": "广东省",
          "City": "广州市",
          "District": "白云区",
          "Address": "白云路 500 号 404 室"
      },
      "OrderLine": [
        {
            "SpuId": 73867287215475272,
            "SkuId": 73867701344275016,
            "Qty": 1
        }
      ]
    }



    grpcurl -plaintext -d '{"Action":"CreateOrder","account_id":60121,"organization_id":60077,"top_organization_id":60066,"ZoneId":666003001,"channel":1,"CouponCode":"PBIXTZBX","PayMethod":1,"Recipient":{"Name":"王永庆","Email":"wolf.sr@metaverse.com","Mobile":"13499192345","Province":"广东省","City":"广州市","District":"白云区","Address":"白云路 500 号 404 室"},"OrderLine":[{"SpuId":73867287215475272,"SkuId":73867701344275016,"Qty":1}]}' localhost:29000 api.mutualaid.Order.CreateOrder

## 下单 -- 单个商品账户余额支付不使用优惠券

    {
      "Action": "CreateOrder",
      "account_id": 60160,
      "organization_id": 1,
      "top_organization_id": 1,
      "ZoneId": 666003001,
      "channel": 1,
      "PayMethod": 3,
      "Recipient": {
          "Name": "柳大庆",
          "Email": "will.sr@metaverse.com",
          "Mobile": "13499192345",
          "Province": "广东省",
          "City": "广州市",
          "District": "白云区",
          "Address": "白云路 500 号 404 室"
      },
      "OrderLine": [
        {
            "SpuId": 73867287215475272,
            "SkuId": 73867701344275016,
            "Qty": 1
        }
      ]
    }



    grpcurl -plaintext -d '{"Action":"CreateOrder","account_id":60160,"organization_id":1,"top_organization_id":1,"ZoneId":666003001,"channel":1,"PayMethod":3,"Recipient":{"Name":"柳大庆","Email":"tinker2022@163.com","Mobile":"13499192345","Province":"广东省","City":"广州市","District":"白云区","Address":"白云路 500 号 404 室"},"OrderLine":[{"SpuId":73867287215475272,"SkuId":73867701344275016,"Qty":1}]}' localhost:29000 api.mutualaid.Order.CreateOrder


## 下单 -- 多个商品不使用优惠券

    {
      "Action": "CreateOrder",
      "account_id": 60160,
      "organization_id": 1,
      "top_organization_id": 1,
      "ZoneId": 666003001,
      "channel": 1,
      "PayMethod": 1,
      "Recipient": {
          "Name": "王永庆",
          "Email": "wolf.sr@metaverse.com",
          "Mobile": "13499192345",
          "Province": "广东省",
          "City": "广州市",
          "District": "白云区",
          "Address": "白云路 500 号 404 室"
      },
      "OrderLine": [
        {
            "SpuId": 73867287215475272,
            "SkuId": 73867701344275016,
            "Qty": 1
        },
        {
            "SpuId": 73868147114904136,
            "SkuId": 73868300190222920,
            "Qty": 2
        }
      ]
    }


    grpcurl -plaintext -d '{"Action":"CreateOrder","account_id":60160,"organization_id":1,"top_organization_id":1,"ZoneId":666003001,"channel":1,"PayMethod":1,"Recipient":{"Name":"王永庆","Email":"wolf.sr@metaverse.com","Mobile":"13499192345","Province":"广东省","City":"广州市","District":"白云区","Address":"白云路 500 号 404 室"},"OrderLine":[{"SpuId":73867287215475272,"SkuId":73867701344275016,"Qty":1},{"SpuId":73868147114904136,"SkuId":73868300190222920,"Qty":2}]}' localhost:29000 api.mutualaid.Order.CreateOrder


## 下单 -- 单个商品不使用优惠券不设置可用区

    grpcurl -plaintext -d '{"Action":"CreateOrder","account_id":60160,"organization_id":1,"top_organization_id":1,"channel":2,"PayMethod":1,"Recipient":{"Name":"王永庆","Email":"wolf.sr@metaverse.com","Mobile":"13499192345","Province":"广东省","City":"广州市","District":"白云区","Address":"白云路 500 号 404 室"},"OrderLine":[{"SpuId":73867287215475272,"SkuId":73867701344275016,"Qty":1},{"SpuId":73868147114904136,"SkuId":73868300190222920,"Qty":2}]}' localhost:29000 api.mutualaid.Order.CreateOrder


## 下单 -- 多个商品不使用优惠券账户余额支付

    {
      "Action": "CreateOrder",
      "account_id": 60160,
      "organization_id": 1,
      "top_organization_id": 1,
      "ZoneId": 666003001,
      "channel": 1,
      "PayMethod": 3,
      "Recipient": {
          "Name": "王永庆",
          "Email": "wolf.sr@metaverse.com",
          "Mobile": "13499192345",
          "Province": "广东省",
          "City": "广州市",
          "District": "白云区",
          "Address": "白云路 500 号 404 室"
      },
      "OrderLine": [
        {
            "SpuId": 73867287215475272,
            "SkuId": 73867701344275016,
            "Qty": 1
        },
        {
            "SpuId": 73868147114904136,
            "SkuId": 73868300190222920,
            "Qty": 2
        }
      ]
    }


    grpcurl -plaintext -d '{"Action":"CreateOrder","account_id":60160,"organization_id":1,"top_organization_id":1,"ZoneId":666003001,"channel":1,"PayMethod":3,"Recipient":{"Name":"王永庆","Email":"wolf.sr@metaverse.com","Mobile":"13499192345","Province":"广东省","City":"广州市","District":"白云区","Address":"白云路 500 号 404 室"},"OrderLine":[{"SpuId":73867287215475272,"SkuId":73867701344275016,"Qty":3},{"SpuId":73868147114904136,"SkuId":73868300190222920,"Qty":2}]}' localhost:29000 api.mutualaid.Order.CreateOrder


## 查订单状态

    {
      "Action": "GetOrderStatus",
      "top_organization_id": 1,
      "OrderNo": "74916676793057085"
    }

    grpcurl -plaintext -d '{ "Action": "GetOrderStatus", "top_organization_id": 1, "OrderNo": "74916676793057085" }' localhost:29000 api.mutualaid.Order.GetOrderStatus

#  用户服务
**需要安装httpie**

## 发送验证码

```shell
echo '{"Phone":"13564130328"}'|http localhost:8080/authcode
```

## 验证码登录

```shell
echo '{"Phone":"13564130328", "AuthCode":"478211"}'|http localhost:8080/auth
```

## 激活用户

```shell
cat<<EOF | http https://mutualaid.ucloud.cn/api/activeuser

{
    "login_code": "Ro1Pi9nIpfXb1eWqZlZJ1o",
    "phone_code": "上海市杨浦区",
    "name": "lance.wang",
    "icon": "https://s.cn.bing.net/th?id=ODLSF.7e8df545-6ead-4949-9ce0-1862720761a5&w=16&h=16&o=6&pid=1.2"
}
EOF
```

## 读取用户信息

```shell
http GET localhost:28000/user Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

## 辅助前端进行签名

```shell
cat<<EOF | http localhost:28000/js_api_sign

{
  "api_url":"http://cn.bing.com"
}
EOF
```
## 我的帮助

```shell
http GET 'localhost:28000/user/aid/offered?page_number=0&page_size=10' Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

## 我的求助

```shell
http GET 'localhost:28000/user/aid/needs?page_number=0&page_size=10' Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

## 求助详情

```shell
http GET 'localhost:28000/user/aid/93768179259539458?id=93768179259539458&latitude=10&longitude=10' Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

## 取电话号码

```shell
http GET 'localhost:28000/wx_phone_number?phone_code=b42e8a2b11d7ac22e9a0093f99f36279416b8f4ac3bb7f1d7b5263fe08dcc57b' Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

## 发现页

```shell
http GET 'https://mutualaid.ucloud.cn/api/discovery?latitude=31.27495384&longitude=121.48361206&page_number=1&page_size=10' Authorization:"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTI2OTE2NTUsImlzcyI6Im11dHVhbGFpZC1jc3IudWNsb3VkLmNuIiwidWlkIjoxMjE5MzExNzU2MjIwODI1OCwib3BlbmlkIjoib2RIUlE0eGFMd21UNDNOdF9JS2RYTExmX3hEYyJ9.vSZsiqA2w27fUbx4OKbhX-89miTPQJ6f_ugUBSm1EuI"
```

# 查询wx api rid详情
```shell

echo '{"rid":"62600b5e-066ae055-2ff7f1b0"}' | http "https://api.weixin.qq.com/cgi-bin/openapi/rid/get?access_token=56_ZweYZnEPeh2wDCn38eEqr-L5rpuUCAw6odDH-BJoqoj8_LxAhIw-KIPOUWhIfwbnGwSM3OMYjOuFezHZgflt2Ue8AFkEKhaSbciBKFzUQhrMKJMCrMMmNdoriuXJCBdZxStYu5zg7BNt55Y1FVPaAIABRL"
```

# 公众号授权登录
```shell

 http "http://localhost:28000/wxoauth2?code=123"
```