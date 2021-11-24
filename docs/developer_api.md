# 开发者API手册
  * [集群模板API](#----api)
    + [1. 创建集群](#1-----)
    + [2. 获取集群列表](#2-------)
    + [3. 创建VPC](#3---vpc)
    + [4. 查看VPC](#4---vpc)
    + [5. 创建子网](#5-----)
    + [6. 查看子网](#6-----)
    + [7. 创建安全组](#7------)
    + [8. 查看安全组](#8------)
    + [9. 创建网络配置](#9-------)
    + [10. 查看region列表](#10---region--)
    + [11. 查看zone列表](#11---zone--)
    + [12. 查看机型列表](#12-------)
    + [13. 获取镜像列表](#13-------)
  * [扩缩容任务API](#-----api)
    + [1. 创建扩容任务](#1-------)
    + [2. 创建缩容任务](#2-------)
    + [3. 查看任务列表](#3-------)
  * [机器API](#--api)
    + [1. 机器列表](#1-----)
    + [2. 机器详情](#2-----)
    + [3. 获取机器数量](#3-------)
  * [费用API](#--api)
    + [1. 单日使用机器总时长](#1----------)
    + [2. 单日使用机器时长明细](#2-----------)

## 集群模板API
### 1. 创建集群
创建用户需要的集群模板<br>

**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/cluster/create </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>name</td>
    <td>string</td>
    <td>是</td>
    <td>集群名称</td>
    <td>test_cluster</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>string</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>string</td>
    <td>是</td>
    <td>区域</td>
    <td>cn-beijing</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>string</td>
    <td>是</td>
    <td>可用区</td>
    <td>cn-beijing-h</td>
  </tr>
  <tr>
    <td>network_config</td>
    <td>object{}</td>
    <td>是</td>
    <td>网络配置信息</td>
    <td>{}</td>
  </tr>
  <tr>
    <td>instance_type</td>
    <td>string</td>
    <td>是</td>
    <td>实例规格</td>
    <td>ecs.s6-c1m1.small</td>
  </tr>
  <tr>
    <td>charge_type</td>
    <td>string</td>
    <td>是</td>
    <td>付费类型：<br>
PostPaid按量付费 <br>
PrePaid包年包月</td>
    <td>PostPaid</td>
  </tr>
  <tr>
    <td>image</td>
    <td>string</td>
    <td>是</td>
    <td>系统镜像id</td>
    <td>m-2ze14bof6m3aadve22aq</td>
  </tr>
  <tr>
    <td>disks</td>
    <td>object{}</td>
    <td>是</td>
    <td>存储配置信息</td>
    <td>{}</td>
  </tr>
  <tr>
    <td>password</td>
    <td>string</td>
    <td>是</td>
    <td>默认实例密码</td>
    <td>ASDqwe123</td>
  </tr>
  <tr>
    <td>account_key</td>
    <td>string</td>
    <td>是</td>
    <td>云账户ak</td>
    <td>LTAI5tAwAMpXAQ78pePcRb6t</td>
  </tr>
  <tr>
    <td>desc</td>
    <td>string</td>
    <td>否</td>
    <td>集群描述</td>
    <td>用于应对突发业务</td>
  </tr>
    <tr>
    <td>tags</td>
    <td>string</td>
    <td>否</td>
    <td>集群标签</td>
    <td>null</td>
</table>


**network_config中的内容**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc</td>
    <td>string</td>
    <td>是</td>
    <td>虚拟私有网络vpc的id</td>
    <td>vpc-2zelmmlfd5c5duibc2xb2</td>
  </tr>
  <tr>
    <td>subnet_id</td>
    <td>string</td>
    <td>是</td>
    <td>子网id</td>
    <td>vsw-2zennaxawzq6sa2fdj8l5</td>
  </tr>
  <tr>
    <td>security_group</td>
    <td>string</td>
    <td>是</td>
    <td>安全组id</td>
    <td>sg-2zefbt9tw0yo1r7vc3ac</td>
  </tr>
  <tr>
    <td>internet_charge_type</td>
    <td>string</td>
    <td>是(需要公网带宽的时候为必传参数)</td>
    <td>网络计费类型。取值范围：<br>
PayByBandwidth：按固定带宽计费。<br>
PayByTraffic（默认）：按使用流量计费。</td>
    <td>PayByTraffic</td>
  </tr>
  <tr>
    <td>internet_charge_type</td>
    <td>string</td>
    <td>是(需要公网带宽的时候为必传参数)</td>
    <td>网络最大带宽(M)</td>
    <td>10</td>
  </tr>
  
</table>


**disks中的内容**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>system_disk</td>
    <td>object</td>
    <td>是</td>
    <td>系统盘配置</td>
    <td></td>
  </tr>
  <tr>
    <td>data_disk</td>
    <td>object</td>
    <td>是</td>
    <td>数据盘配置</td>
    <td></td>
  </tr>
</table>

**system_disk和data_disk中的内容**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>category</td>
    <td>string</td>
    <td>是</td>
    <td>系统盘(数据盘)类型</td>
    <td>cloud_efficiency</td>
  </tr>
  <tr>
    <td>size</td>
    <td>int</td>
    <td>是</td>
    <td>系统盘(数据盘)大小(G)</td>
    <td>40</td>
  </tr>
</table>


**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**请求示例**
```JSON
{
    "name":"测试集群",
    "provider":"AlibabaCloud",
    "account_key":"LTAI5t7qCv**Fh3hzSYpSv",
    "charge_type":"PostPaid",
    "region_id":"cn-qingdao",
    "zone_id":"cn-qingdao-b",
    "instance_type":"ecs.n1.tiny",
    "image":"centos_7_6_x64_20G_alibase_20211030.vhd",
    "password":"********",
    "desc":"测试用的集群",
    "network_config":{
        "vpc":"vpc-m5***6tgd",
        "subnet_id":"vsw-m5***4y3xs6ivwc",
        "security_group":"sg-m5***2cu2f5x8"
    },
    "disks":{
        "system_disk":{
            "category":"cloud_efficiency",
            "size":50,
        },
        "data_disk":[
            {
                "category":"cloud_efficiency",
                "size":50,
            }
        ]
    }
}
```

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "msg":"success",
    "data":"测试集群"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**
<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>



### 2. 获取集群列表
获取本账户下所有的集群信息<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/cluster/describe_all </td>
  </tr>
</table>


**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>String</td>
    <td>否</td>
    <td>云账户(不传默认查询组织下所有集群)</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>子属性</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>cluster_list[]</td>
    <td></td>
    <td>Array</td>
    <td>是</td>
    <td>集群列表</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>cluster_name</td>
    <td>String</td>
    <td>是</td>
    <td>集群名</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>AlibabaCloud</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>account</td>
    <td>String</td>
    <td>是</td>
    <td>云账号</td>
    <td>aagjege</td>
  </tr>
  <tr>
    <td></td>
    <td>total_remainder</td>
    <td>String</td>
    <td>是</td>
    <td>配额用量/余量</td>
    <td>8/200</td>
  </tr>
  <tr>
    <td></td>
    <td>tcreate_at</td>
    <td>String</td>
    <td>是</td>
    <td>创建时间</td>
    <td>2021-11-03 17:01:44</td>
  </tr>
  <tr>
    <td>pager</td>
    <td>Pager</td>
    <td>String</td>
    <td>是</td>
    <td>分页参数</td>
    <td></td>
  </tr>
</table>


**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>


**响应示例**

正常返回结果：<br>
```JSON
{
  "code": 200,
  "data": {
    "cluster_list": [
      {
        "cluster_id": "1319",
        "cluster_name": "gf.bridgx.online",
        "provider": "AlibabaCloud",
        "account": "LTAI5tAwAMpXAQ78pePcRb6t",
        "create_at": "2021-11-02 06:09:38 +0800 CST",
        "create_by": ""
      },
      {
        "cluster_id": "1332",
        "cluster_name": "test6",
        "provider": "AlibabaCloud",
        "account": "LTAI5tAwAMpXAQ78pePcRb6t",
        "create_at": "2021-11-04 14:22:15 +0800 CST",
        "create_by": ""
      }
    ],
    "pager": {
      "page_number": 1,
      "page_size": 10,
      "total": 5
    }
  },
  "msg": "success"
}
```

异常返回结果：

```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
</table>


### 3. 创建VPC
一个VPC只能指定一个网段，网段的范围包括10.0.0.0/8、172.16.0.0/12和192.168.0.0/16以及它们的子网，网段的掩码为8~24位，默认为172.16.0.0/12。
VPC创建后无法修改网段。<br>
每个VPC支持云资源使用的私网网络地址数量为60,000个，且无法提升配额。<br>
创建VPC后，会自动创建一个路由器和一个路由表。<br>
每个VPC支持三个用户侧网段。如果多个用户侧网段之间存在包含关系，掩码较短的网段实际生效。例如10.0.0.0/8和10.1.0.0/16中，10.0.0.0/8实际生效。<br>

**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/vpc/create</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>是</td>
    <td>目标地域ID</td>
    <td>cn-hangzhou</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>String</td>
    <td>否</td>
    <td>VPC的网段。您可以使用以下网段或其子集作为VPC的网段：<br>
172.16.0.0/12（默认值）<br>
10.0.0.0/8<br>
192.168.0.0/16。</td>
    <td>172.16.0.0/12</td>
  </tr>  
  <tr>
    <td>vpc_name</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的名称:
长度为2～64个字符，必须以字母或中文开头，<br>
可包含数字、半角句号（.）、下划线（_）和短划线（-），<br>
但不能以http:// 或https:// 开头。</td>
    <td>abc</td>
  </tr>
  <tr>
    <td>ak</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商账户下的 AK</td>
    <td>dasdadasdasd</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>string</td>
    <td>是</td>
    <td>创建成功后的vpcid</td>
    <td>vpc-asd**asdasda</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-hangzhou</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>172.16.0.0/12</td>
  </tr>
  <tr>
    <td>vpc_name</td>
    <td>abc</td>
  </tr>
  <tr>
    <td>ak</td>
    <td>dasdadasdasd</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "msg":"success",
    "data":"vsw-m5evh*****s6ivwc"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 4. 查看VPC
根据VPC的名字查找用户已经创建的vpc信息

**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/vpc/describe</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>否</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>否</td>
    <td>VPC所在的地域ID</td>
    <td>cn-hangzhou</td>
  </tr>
  <tr>
    <td>vpc_name</td>
    <td>String</td>
    <td>否</td>
    <td>VPC的名称</td>
    <td>vpc-1</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>


**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>create_at</td>
    <td>String</td>
    <td>是</td>
    <td>VPC创建时间</td>
    <td>2018-04-18T15:02:37Z</td>
  </tr>
  <tr>
    <td>status</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的状态，取值：<br>
Pending：配置中。<br>
Available：可用。</td>
    <td>Available</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的ID</td>
    <td>vpc-bp1qpo0kug3a20qqe****</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的IPv4网段</td>
    <td>192.168.0.0/16</td>
  </tr>
  <tr>
    <td>switch_ids</td>
    <td>String []</td>
    <td>是</td>
    <td>VPC中交换机的列表</td>
    <td>vsw-bp1nhbnpv2blyz8dl****</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String []</td>
    <td>是</td>
    <td>所属的云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>vpc_name</td>
    <td>String []</td>
    <td>否</td>
    <td>VPC的名称</td>
    <td>vpc-1</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-qingdao</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>50</td>
  </tr>
</table>


**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "Vpcs":[
            {
                "VpcId":"vpc-m5**swmv796tgd",
                "VpcName":"vpc测试一",
                "CidrBlock":"",
                "SwitchIds":"",
                "Provider":"AlibabaCloud",
                "Status":"",
                "CreateAt":"2021-11-11 11:13:34 +0800 CST"
            },
            {
                "VpcId":"vpc-m5e50cwcefjgxvbjs1ud5",
                "VpcName":"vpc测试二",
                "CidrBlock":"",
                "SwitchIds":"",
                "Provider":"AlibabaCloud",
                "Status":"",
                "CreateAt":"2021-11-05 11:40:16 +0800 CST"
            }
        ],
        "Pager":{
            "PageNumber":1,
            "PageSize":50,
            "Total":2
        }
    },
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**
<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>



### 5. 创建子网
- 每个VPC内的交换机数量不能超过150个。
- 每个交换机网段的第1个和最后3个IP地址为系统保留地址。例如192.168.1.0/24的系统保留地址为192.168.1.0、192.168.1.253、192.168.1.254和192.168.1.255。
- 交换机下的云产品实例数量不允许超过VPC剩余的可用云产品实例数量（15000减去当前云产品实例数量）。
- 一个云产品实例只能属于一个交换机。
- 交换机不支持组播和广播。
- 交换机创建成功后，无法修改网段

**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/subnet/create</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>String</td>
    <td>是</td>
    <td>可用区ID</td>
    <td>cn-hangzhou-g</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的网段。交换机网段要求如下：<br>
交换机的网段的掩码长度范围为16～29位。<br>
交换机的网段必须从属于所在VPC的网段。<br>
交换机的网段不能与所在VPC中路由条目的目标网段相同，<br>
      但可以是目标网段的子集。</td>
    <td>172.16.0.0/24</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>要创建的交换机所属的VPC ID</td>
    <td>vpc-257gqcdfvx6n****</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的名称,名称长度为2～128个字符，必须以字母或中文开头，<br>
      但不能以http:// 或https:// 开头。</td>
    <td>VSwitch-1</td>
  </tr>
</table>


**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>string</td>
    <td>是</td>
    <td>创建成功的 Switch id</td>
    <td> "asdasdasdas"</td>
  </tr>
</table>


**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>cn-qingdao-b</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>172.16.0.0/24</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>vpc-257gqcdfvx6n****</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>VSwitch-1</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "msg":"success",
    "data":"vsw-m5evh***4y3xs6ivwc"
}
```

异常返回结果    
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 6. 查看子网
来查找子网id信息<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/subnet/describe</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>要查询的交换机所属VPC的ID</td>
    <td>vpc-257gqcdfvx6n****</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>String</td>
    <td>否</td>
    <td>要查询的交换机的名称。(不传默认查询VPC下全部 switch)</td>
    <td>VSwitch-1</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>交换机所属VPC的ID</td>
    <td>vpc-257gqcdfvx6n****</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>String</td>
    <td>是</td>
    <td>要查询的交换机的名称</td>
    <td>VSwitch-1</td>
  </tr>
  <tr>
    <td>status</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的状态，取值：<br>
Pending：配置中。<br>
Available：可用。</td>
    <td>Available</td>
  </tr>
  <tr>
    <td>create_at</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的创建时间</td>
    <td>2018-01-18T12:43:57Z</td>
  </tr>
  <tr>
    <td>is_default</td>
    <td>Boolean</td>
    <td>是</td>
    <td>是否是默认交换机。<br>
true：是默认交换机。<br>
false：非默认交换机。</td>
    <td>true</td>
  </tr>
  <tr>
    <td>available_ip_address_count</td>
    <td>Long</td>
    <td>是</td>
    <td>交换机中可用的IP地址数量</td>
    <td>1</td>
  </tr>
  <tr>
    <td>switch_id</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的ID</td>
    <td>vsw-25bcdxs7pv1****</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的IPv4网段</td>
    <td>172.16.0.0/24</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>String</td>
    <td>是</td>
    <td>可用区ID</td>
    <td>cn-hangzhou-g</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>vpc-257gqcdfvx6n****</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>"第一台交换机"</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "Switches":[
            {
                "VpcId":"vpc-257gqcdfvx6n****",
                "SwitchId":"vsw-m5evh***4y3xs6ivwc",
                "ZoneId":"cn-qingdao-b",
                "SwitchName":"第一台交换机",
                "CidrBlock":"172.16.0.0/24",
                "VStatus":"",
                "CreateAt":"2021-11-02 17:50:37 +0800 CST",
                "IsDefault":"N",
                "AvailableIpAddressCount":0
            }
        ],
        "Pager":{
            "PageNumber":1,
            "PageSize":50,
            "Total":1
        }
    },
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 7. 创建安全组
- 安全组的API文档中，流量的发起端为源端（Source），数据传输的接收端为目的端（Dest）。
- 出方向和入方向安全组规则总和不能超过200条。
- 安全组规则优先级（Priority）可选范围为1~100。数字越小，代表优先级越高。
- 优先级相同的安全组规则，以拒绝访问（drop）的规则优先。
- 源端设备可以是指定的IP地址范围（SourceCidrIp、Ipv6SourceCidrIp、SourcePrefixListId），也可以是其他安全组（SourceGroupId）中的ECS实例。
- 如果匹配的安全组规则已存在，此次AuthorizeSecurityGroup调用成功，但不会增加规则数。

**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/security_group/create_with_rule </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>安全组所属VPC ID</td>
    <td>vpc-bp1opxu1zkhn00gzv****</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>String</td>
    <td>是</td>
    <td>安全组名称。长度为2~128个英文或中文字符。<br>
      必须以大小字母或中文开头，不能以 http:// 和<br>
      https:// 开头。可以包含数字、半角冒号（:）、<br>
      下划线（_）或者连字符（-）。默认值：空。</td>
    <td>testSecurityGroupName</td>
  </tr>
  <tr>
    <td>rules</td>
    <td>[]Rule</td>
    <td>是</td>
    <td>规则</td>
    <td></td>
  </tr>
  <tr>
    <td>rules</td>
    <td>[]Rule</td>
    <td>是</td>
    <td>规则</td>
    <td></td>
  </tr>
</table>



**Rule**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>protocol</td>
    <td>String</td>
    <td>是</td>
    <td>传输层协议。取值大小写敏感。取值范围：<br>
    tcp<br>
    udp<br>
    icmp<br>
    gre<br>
    all：支持所有协议</td>
    <td>tcp</td>
  </tr>
  <tr>
    <td>port_range</td>
    <td>String</td>
    <td>是</td>
    <td>目的端安全组开放的传输层协议相关的端口范围。取值范围：<br>
TCP/UDP协议：取值范围为1~65535。使用斜线（/）隔开起始端口和终止端口。例如：1/200<br>
ICMP协议：-1/-1<br>
GRE协议：-1/-1<br>
IpProtocol取值为all：-1/-1</td>
    <td>22/22</td>
  </tr>
  <tr>
    <td>direction</td>
    <td>String</td>
    <td>是</td>
    <td>安全组规则的方向：取值范围：<br>
egress：出方向<br>
ingress：入方向</td>
    <td>ingress</td>
  </tr>
  <tr>
    <td>group_id</td>
    <td>String</td>
    <td>二选一</td>
    <td>需要设置访问权限的(入或者出)安全组ID</td>
    <td>sg-bp67acfmxazb4p****</td>
  </tr>
  <tr>
    <td>cidr_ip</td>
    <td>String</td>
    <td>二选一</td>
    <td>需要设置访问权限的(入或者出)IPv4 CIDR地址块</td>
    <td>10.0.0.0/8</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>string</td>
    <td>是</td>
    <td>创建成功的安全组 id</td>
    <td> "asdadasdasdad"</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>vpc-bp1opxu1zkhn00gzv****</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>testSecurityGroupName</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
  <tr>
    <td>rules</td>
    <td>
      [{
"protocol":"tcp",<br>
"port_range":"22/22",<br>
"direction":"ingress",<br>
"group_id":"sg-bp67acfmxazb4p****",<br>
"cidr_ip":"10.0.0.0/8"
    }]
    </td>
  </tr>
</table>



**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "msg":"success",
    "data":"vsw-m5evh***4y3xs6ivwc"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>



### 8. 查看安全组
查看已经创建的安全组<br>

**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/security_group/describe </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>否</td>
    <td>安全组所属地域ID</td>
    <td>cn-hangzhou</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>安全组所属VPC ID</td>
    <td>vpc-bp1opxu1zkhn00gzv****</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>String</td>
    <td>否</td>
    <td>安全组名称。(不传默认查 VPC 下所有)</td>
    <td>testSecurityGroupName</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>t1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>create_at</td>
    <td>String</td>
    <td>是</td>
    <td>创建时间。按照ISO8601标准表示，并需要使用UTC时间。<br>
      格式为：yyyy-MM-ddThh:mmZ。</td>
    <td>2021-08-31T03:12:29Z</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>String</td>
    <td>是</td>
    <td>安全组所属的专有网络</td>
    <td>vpc-bp67acfmxazb4p****</td>
  </tr>
  <tr>
    <td>security_group_id</td>
    <td>String</td>
    <td>是</td>
    <td>安全组ID</td>
    <td>sg-bp67acfmxazb4p****</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>String</td>
    <td>是</td>
    <td>安全组名称</td>
    <td>SGTestName</td>
  </tr>
  <tr>
    <td>security_group_type</td>
    <td>String</td>
    <td>是</td>
    <td>安全组类型。可能值：<br>
normal：普通安全组<br>
enterprise：企业安全组</td>
    <td>normal</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-hangzhou</td>
  </tr>
  <tr>
    <td>vpc_id</td>
    <td>vpc-bp1opxu1zkhn00gzv****</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td></td>
  </tr>
  <tr>
    <td>spage_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>


**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "Groups":[
            {
                "VpcId":"vpc-bp1opxu1zkhn00gzv****",
                "SecurityGroupId":"sg-m5ebc***v2cu2f5x8",
                "SecurityGroupName":"测试的第一个安全组",
                "SecurityGroupType":"normal",
                "CreateAt":"2021-11-02 18:03:01 +0800 CST"
            }
        ],
        "Pager":{
            "PageNumber":1,
            "PageSize":15,
            "Total":1
        }
    },
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>




### 9. 创建网络配置
通过一个api即可创建vpc、子网及安全组<br>
**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/network_config/create</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>是</td>
    <td>目标地域ID</td>
    <td>cn-qingdao</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的网段。您可以使用以下网段或其子集作为VPC的网段：<br>
172.16.0.0/12（默认值）。<br>
10.0.0.0/8。<br>
192.168.0.0/16。</td>
    <td>172.16.0.0/12</td>
  </tr>
  <tr>
    <td>vpc_name</td>
    <td>String</td>
    <td>是</td>
    <td>VPC的名称:长度为2～128个字符，必须以字母或中文开头，<br>
      可包含数字、半角句号（.）、下划线（_）和短划线（-），<br>
      但不能以http:// 或https:// 开头。</td>
    <td>abc</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>String</td>
    <td>是</td>
    <td>可用区ID</td>
    <td>cn-qigndao-b</td>
  </tr>
  <tr>
    <td>switch_cidr_block</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的网段。交换机网段要求如下：<br>
交换机的网段的掩码长度范围为16～29位。<br>
交换机的网段必须从属于所在VPC的网段。<br>
交换机的网段不能与所在VPC中路由条目的目标网段相同，但可以是目标网段的子集。
</td>
    <td>172.16.0.0/24</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>String</td>
    <td>是</td>
    <td>交换机的名称:长度为2～128个字符，必须以字母或中文开头，<br>
      但不能以http:// 或https:// 开头。</td>
    <td>VSwitch-1</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>String</td>
    <td>是</td>
    <td>安全组名称。长度为2~128个英文或中文字符。<br>
      必须以大小字母或中文开头，不能以 http:// 和https:// 开头。<br>
      可以包含数字、半角冒号（:）、下划线（_）或者连字符（-）。默认值：空。</td>
    <td>testSecurityGroupName</td>
  </tr>
  <tr>
    <td>security_group_type</td>
    <td>String</td>
    <td>是</td>
    <td>安全组类型。可能值：<br>
normal：普通安全组<br>
enterprise：企业安全组</td>
    <td>normal</td>
  </tr>
  <tr>
    <td>ak</td>
    <td>String</td>
    <td>是</td>
    <td></td>
    <td></td>
  </tr>
</table>


**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>string</td>
    <td>是</td>
    <td>创建成功的vpc_id</td>
    <td>sdasd1qwdasdasd</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-qingdao</td>
  </tr>
  <tr>
    <td>cidr_block</td>
    <td>172.16.0.0/12</td>
  </tr>
  <tr>
    <td>vpc_name</td>
    <td>abc</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>cn-qigndao-b</td>
  </tr>
  <tr>
    <td>switch_cidr_block</td>
    <td>172.16.0.0/24</td>
  </tr>
  <tr>
    <td>switch_name</td>
    <td>VSwitch-1</td>
  </tr>
  <tr>
    <td>security_group_name</td>
    <td>testSecurityGroupName</td>
  </tr>
  <tr>
    <td>security_group_type</td>
    <td>normal</td>
  </tr>
  <tr>
    <td>ak</td>
    <td>asdadasdsadas</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "msg":"success",
    "data":"vsw-m5evh*****s6ivwc"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 10. 查看region列表
查看云厂商的地域列表<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/region/list </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
</table>
 
**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":[
        {
            "RegionId":"cn-qingdao",
            "LocalName":"华北 1"
        },
        {
            "RegionId":"cn-beijing",
            "LocalName":"华北 2"
        }
    ],
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data": null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 11. 查看zone列表
通过云厂商的地域查看该地域下的可用区列表<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/zone/list</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>是</td>
    <td>地域ID</td>
    <td> cn-qingdao</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-qingdao</td>
  </tr>
</table>


**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":[
        {
            "ZoneId":"cn-qingdao-b",
            "LocalName":"华北 1 可用区 B"
        },
        {
            "ZoneId":"cn-qingdao-c",
            "LocalName":"华北 1 可用区 C"
        }
    ],
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 12. 查看机型列表
通过云厂商的地域和可用区查看该地域下及可用区下的机型列表<br>

**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance_type/list</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>是</td>
    <td>地域ID</td>
    <td>cn-qingdao</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>String</td>
    <td>是</td>
    <td>可用区</td>
    <td>cn-qingdao-b</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_type_family</td>
    <td>String</td>
    <td>是</td>
    <td>机型规格族</td>
    <td>ecs.g6</td>
  </tr>
  <tr>
    <td>instance_type</td>
    <td>String</td>
    <td>是</td>
    <td>机型名称</td>
    <td>ecs.g6.large</td>
  </tr>
  <tr>
    <td>core</td>
    <td>int</td>
    <td>是</td>
    <td>cpu核数</td>
    <td>4</td>
  </tr>
  <tr>
    <td>memory</td>
    <td>int</td>
    <td>是</td>
    <td>内存大小，单位G</td>
    <td>8</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-qingdao</td>
  </tr>
  <tr>
    <td>zone_id</td>
    <td>cn-qingdao-b</td>
  </tr>
</table>


**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":[
        {
            "instance_type_family":"ecs.e3",
            "instance_type":"ecs.e3.large",
            "core":4,
            "memory":32
        },
        {
            "instance_type_family":"ecs.r6",
            "instance_type":"ecs.r6.13xlarge",
            "core":52,
            "memory":384
        }
    ],
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>




### 13. 获取镜像列表
通过云厂商的地域和可用区查看该地域下及可用区下的机型列表<br>
**请求地址**
<table>
  <tr>
    <td>GET 方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/image/list</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>String</td>
    <td>是</td>
    <td>区域</td>
    <td>cn-qingdao</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>


**data重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>os_type</td>
    <td>String</td>
    <td>是</td>
    <td>操作系统类型</td>
    <td>linux</td>
  </tr>
  <tr>
    <td>os_name</td>
    <td>String</td>
    <td>是</td>
    <td>操作系统名称</td>
    <td>Windows Server 2016 数据中心版 64位中文版</td>
  </tr>
  <tr>
    <td>image_id</td>
    <td>String</td>
    <td>是</td>
    <td>镜像ID</td>
    <td>m-bp1g7004ksh0oeuc****</td>
  </tr>
</table>


**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td>cn-qingdao</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":[
        {
            "OsType":"linux",
            "OsName":"CentOS  7.6 64位",
            "ImageId":"centos_7_6_x64_20G_alibase_20211030.vhd"
        },
        {
            "OsType":"linux",
            "OsName":"Gentoo  13  64bit",
            "ImageId":"gentoo13_64_40G_aliaegis_20160222.vhd"
        }
    ],
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>



## 扩缩容任务API
### 1. 创建扩容任务
扩大某集群的机器数量<br>

**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/cluster/expand </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>task_name</td>
    <td>String</td>
    <td>是</td>
    <td>任务名称</td>
    <td>expand_task</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>是</td>
    <td>集群的名称</td>
    <td>gf.metrics.test</td>
  </tr>
  <tr>
    <td>count</td>
    <td>Int</td>
    <td>是</td>
    <td>扩容的机器数量</td>
    <td>10</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>


**请求示例**
```JSON
{
    "cluster_name":"gf.bridgx.online",
    "task_name":"aaa",
    "count":1
}
```
**响应示例**

正常返回结果：
```JSON
{
  "code": 200,
  "data": 69762**4493870,
  "msg": "success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>




### 2. 创建缩容任务
缩小某集群的机器数量，如果指定了IP会按照指定IP进行缩容，不指定IP会随机选择count台机器进行缩容。<br>
**请求地址**
<table>
  <tr>
    <td>POST方法</td>
  </tr>
  <tr>
    <td>POST /api/v1/cluster/expand </td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>task_name</td>
    <td>String</td>
    <td>是</td>
    <td>任务名称</td>
    <td>expand_task</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>是</td>
    <td>集群的名称</td>
    <td>gf.metrics.test</td>
  </tr>
  <tr>
    <td>ips</td>
    <td>String</td>
    <td>否</td>
    <td>缩容的ip地址</td>
    <td>["10.192.220.195", "10.192.220.196", "10.192.220.197"]</td>
  </tr>
  <tr>
    <td>count</td>
    <td>Int</td>
    <td>是</td>
    <td>扩容的机器数量</td>
    <td>10</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**请求示例**
```JSON
{
    "cluster_name":"gf.bridgx.online",
    "task_name":"asas",
    "count":1
}
```
**响应示例**

正常返回结果：
```JSON
{
  "code": 200,
  "data": 69762**4493871,
  "msg": "success"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data": null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>



### 3. 查看任务列表
查询账号下的任务列表，不传account默认查询全部，传了account查询特定account下的。<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/task/list</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>String</td>
    <td>否</td>
    <td>云账户</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data中的重要内容**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>task_id</td>
    <td>String</td>
    <td>是</td>
    <td>任务ID</td>
    <td>4803680234646135</td>
  </tr>
  <tr>
    <td>task_name</td>
    <td>String</td>
    <td>是</td>
    <td>任务名</td>
    <td>xx紧急扩容xx</td>
  </tr>
  <tr>
    <td>task_action</td>
    <td>String</td>
    <td>是</td>
    <td>任务动作类型</td>
    <td>扩容、缩容</td>
  </tr>
  <tr>
    <td>status</td>
    <td>String</td>
    <td>是</td>
    <td>任务状态</td>
    <td>待执行、执行中、执行完毕</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>是</td>
    <td>集群名称</td>
    <td>gf.xx.aa</td>
  </tr>
  <tr>
    <td>create_at</td>
    <td>String</td>
    <td>是</td>
    <td>创建时间</td>
    <td>2021-11-03 16:54:46</td>
  </tr>
  <tr>
    <td>execute_time</td>
    <td>int</td>
    <td>是</td>
    <td>执行时间（单位：秒）</td>
    <td>18</td>
  </tr>
  <tr>
    <td>finish_at</td>
    <td>string</td>
    <td>是</td>
    <td>完成时间</td>
    <td>2021-11-03 16:55:36</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
  "code": 200,
  "data": {
    "task_list": [
      {
        "task_id": "6978292825513784",
        "task_name": "",
        "task_action": "EXPAND",
        "status": "SUCCESS",
        "cluster_name": "gf.scheduler.test",
        "create_at": "2021-11-18 11:23:07 +0800 CST",
        "execute_time": 11,
        "finish_at": "2021-11-18 11:23:18 +0800 CST"
      },
      {
        "task_id": "6692774019653071",
        "task_name": "test再次缩容",
        "task_action": "SHRINK",
        "status": "SUCCESS",
        "cluster_name": "gf.bridgxine",
        "create_at": "2021-11-16 12:06:44 +0800 CST",
        "execute_time": 4,
        "finish_at": "2021-11-16 12:06:48 +0800 CST"
      },
      {
        "task_id": "6692475385208271",
        "task_name": "test",
        "task_action": "SHRINK",
        "status": "SUCCESS",
        "cluster_name": "gf.bridgx.online",
        "create_at": "2021-11-16 12:03:46 +0800 CST",
        "execute_time": 5,
        "finish_at": "2021-11-16 12:03:51 +0800 CST"
      }
    ],
    "pager": {
      "page_number": 1,
      "page_size": 50,
      "total": 52
    }
  },
  "msg": "success"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


## 机器API
### 1. 机器列表
获取本账户下所有的机器信息<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance/describe_all</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>String</td>
    <td>否</td>
    <td>云账户</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>instance_id</td>
    <td>String</td>
    <td>否</td>
    <td>实例ID，精确匹配</td>
    <td>i-xaf23fasdc1edg</td>
  </tr>
  <tr>
    <td>ip</td>
    <td>String</td>
    <td>否</td>
    <td>内网或外网IP，精确匹配</td>
    <td>10.12.13.1</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>String</td>
    <td>否</td>
    <td>特定云厂商，精确匹配</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>否</td>
    <td>集群名称</td>
    <td>集群一</td>
  </tr>
  <tr>
    <td>status</td>
    <td>String</td>
    <td>否</td>
    <td>状态</td>
    <td>running,deleted</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>


**data中的重要内容**
<table>
  <tr>
    <td>名称</td>
    <td>子属性</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_list</td>
    <td>0</td>
    <td>[]</td>
    <td>是</td>
    <td>机器列表</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>instance_id</td>
    <td>String</td>
    <td>是</td>
    <td>实例id</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>ip_inner</td>
    <td>String</td>
    <td>是</td>
    <td>内网IP</td>
    <td>10.208.28.126[内网]</td>
  </tr>
  <tr>
    <td></td>
    <td>ip_outer</td>
    <td>String</td>
    <td>否</td>
    <td>外网IP</td>
    <td>106.14.169.121[公网]</td>
  </tr>
  <tr>
    <td></td>
    <td>provider</td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td></td>
    <td>cluster_name</td>
    <td>String</td>
    <td>是</td>
    <td>所属集群</td>
    <td>gf.bridgx.online</td>
  </tr>
  <tr>
    <td></td>
    <td>instance_type</td>
    <td>String</td>
    <td>是</td>
    <td>机型</td>
    <td>ecs.7c.large</td>
  </tr>
  <tr>
    <td></td>
    <td>create_at</td>
    <td>String</td>
    <td>是</td>
    <td>创建时间</td>
    <td>2021.10.29  18：25：13</td>
  </tr>
  <tr>
    <td></td>
    <td>status</td>
    <td>String</td>
    <td>是</td>
    <td>机器状态</td>
    <td>初始化完成</td>
  </tr>
  <tr>
    <td>pager</td>
    <td>Pager</td>
    <td>String</td>
    <td>是</td>
    <td>是</td>
    <td>分页参数</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>instance_id</td>
    <td>i-xaf23fasdc1edg</td>
  </tr>
  <tr>
    <td>ip</td>
    <td>10.12.13.1</td>
  </tr>
  <tr>
    <td>provider</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>集群一</td>
  </tr>
  <tr>
    <td>status</td>
    <td>running,deleted</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "instance_list":[
            {
                "instance_id":"i-2ze40**rrjk7mi6",
                "ip_inner":"10.192.221.25",
                "ip_outer":"",
                "provider":"AlibabaCloud",
                "create_at":"2021-11-12 09:38:31 +0800 CST",
                "status":"Deleted",
                "startup_time":0,
                "cluster_name":"gf.bridgx.online",
                "instance_type":"ecs.s6-c1m1.small"
            },
            {
                "instance_id":"i-2ze25xv**vu06m0p2",
                "ip_inner":"10.192.221.123",
                "ip_outer":"",
                "provider":"AlibabaCloud",
                "create_at":"2021-11-12 19:58:21 +0800 CST",
                "status":"Deleted",
                "startup_time":5,
                "cluster_name":"gf.bridgx.online",
                "instance_type":"ecs.s6-c1m1.small"
            }
        ],
        "pager":{
            "page_number":1,
            "page_size":10,
            "total":1747
        }
    },
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 2. 机器详情
获取某机器的详细信息<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance/id/describe</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_id</td>
    <td>String</td>
    <td>是</td>
    <td>机器的id</td>
    <td>i-2ze40hb376ihrrjk7mi6</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td>{}</td>
  </tr>
</table>


**data中的重要内容**

<table>
  <tr>
    <td>名称</td>
    <td>子属性</td>
    <td>孙属性</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_id</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>机器的id</td>
    <td></td>
  </tr>
  <tr>
    <td>provider</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>云厂商</td>
    <td>AlibabaCloud</td>
  </tr>
  <tr>
    <td>region_id</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>可用区</td>
    <td>cn-beijing-h</td>
  </tr>
  <tr>
    <td>create_at</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>创建时间</td>
    <td>2021-10-29 16：23：24</td>
  </tr>
  <tr>
    <td>image_id</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>镜像id</td>
    <td></td>
  </tr>
  <tr>
    <td>instance_type</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>实例规格</td>
    <td>4核16G</td>
  </tr>
  <tr>
    <td>storage_config</td>
    <td></td>
    <td></td>
    <td>{object}</td>
    <td>是</td>
    <td>存储配置</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>system_disk_type</td>
    <td>String</td>
    <td>是</td>
    <td>系统盘类型</td>
    <td>cloud_efficiency</td>
  </tr>
  <tr>
    <td></td>
    <td>system_disk_size</td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>系统盘大小</td>
    <td>40G</td>
  </tr>
  <tr>
    <td></td>
    <td>data_disks</td>
    <td></td>
    <td>[]</td>
    <td>否</td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td>data_disk_type</td>
    <td>String </td>
    <td>否</td>
    <td>数据盘类型</td>
    <td>cloud_efficiency</td>
  </tr>
  <tr>
    <td></td>
    <td></td>
    <td>data_disk_size</td>
    <td>String </td>
    <td>否</td>
    <td>数据盘大小</td>
    <td>40G</td>
  </tr>
  <tr>
    <td></td>
    <td>data_disk_num</td>
    <td></td>
    <td>int</td>
    <td>是</td>
    <td>数据盘个数</td>
    <td>4</td>
  </tr>
  <tr>
    <td>network_config</td>
    <td></td>
    <td></td>
    <td>{object}</td>
    <td>是</td>
    <td>网络配置</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>vpc_name</td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>VPC的名称</td>
    <td>testvpc</td>
  </tr>
  <tr>
    <td></td>
    <td>subnet_id_name</td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>子网名字</td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>security_group_name</td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>安全组名字</td>
    <td></td>
  </tr>
  <tr>
    <td>ip_outer</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>公网ip</td>
    <td></td>
  </tr>
  <tr>
    <td>ip_inner</td>
    <td></td>
    <td></td>
    <td>String</td>
    <td>是</td>
    <td>内网ip</td>
    <td></td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_id</td>
    <td>i-2ze40hb**hrrjk7mi6</td>
  </tr>
</table>


**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "instance_id":"i-2ze40hb**hrrjk7mi6",
        "provider":"AlibabaCloud",
        "region_id":"cn-beijing",
        "image_id":"m-2ze**m3aadve22aq",
        "instance_type":"ecs.s6-c1m1.small",
        "ip_inner":"10.192.221.25",
        "ip_outer":"",
        "create_at":"2021-11-12 09:38:31 +0800 CST",
        "storage_config":{
            "system_disk_type":"cloud_efficiency",
            "system_disk_size":40,
            "data_disks":[
                {
                    "data_disk_type":"cloud_efficiency",
                    "data_disk_size":100
                }
            ],
            "data_disk_num":1
        },
        "network_config":{
            "vpc_name":"vpc-2zelmmlf**c2xb2",
            "subnet_id_name":"vsw-2ze**q6sa2fdj8l5",
            "security_group_name":"sg-2zefbt9tw0y***7vc3ac"
        }
    },
    "msg":"success"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 3. 获取机器数量
获取本账户下运行的机器数量<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance/num</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>String</td>
    <td>否</td>
    <td>云账户</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>否</td>
    <td>集群名称</td>
    <td>测试集群</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data 中的重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_num</td>
    <td>int64</td>
    <td>是</td>
    <td>实例数量</td>
    <td>2</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>account</td>
    <td>LTAI5tAWAM</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>测试集群</td>
  </tr>
</table>

**响应示例**

正常返回结果：
```JSON
{
    "code":200,
    "data":{
        "instance_num":1
    },
    "msg":"success"
}
```

异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>




## 费用API
### 1. 单日使用机器总时长
指定集群，返回特定集群的使用时长，否则返回当前账号下关联全部集群的总时长。<br>
**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance/usage_total</td>
  </tr>
</table>

**请求参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>否</td>
    <td>集群名</td>
    <td>gf.bridgx.online</td>
  </tr>
  <tr>
    <td>date</td>
    <td>String</td>
    <td>是</td>
    <td>yyyy-dd-mm</td>
    <td>2021-10-11</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>


**data 中的重要参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>usage_total</td>
    <td>int</td>
    <td>是</td>
    <td>使用时长单位秒</td>
    <td>1800</td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>gf.bridgx.online</td>
  </tr>
  <tr>
    <td>date</td>
    <td>2021-10-11</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>

**响应示例**

正常返回结果:
```JSON
{
    "code": 200,
    "msg": "success",
    "data": {
        "usage_total": 1800, // 表示1800s
    }
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```

**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table>


### 2. 单日使用机器时长明细 

**请求地址**
<table>
  <tr>
    <td>GET方法</td>
  </tr>
  <tr>
    <td>GET /api/v1/instance/usage_statistics</td>
  </tr>
</table>

**请求参数**

<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>String</td>
    <td>否</td>
    <td>集群名</td>
    <td>gf.bridgx.online</td>
  </tr>
  <tr>
    <td>date</td>
    <td>String</td>
    <td>是</td>
    <td>yyyy-dd-mm</td>
    <td>2021-10-11</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>int32</td>
    <td>否</td>
    <td>默认第一页</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>int32</td>
    <td>否</td>
    <td>默认 10 最大 50</td>
    <td>15</td>
  </tr>
</table>

**返回参数**
<table>
  <tr>
    <td>名称</td>
    <td>类型</td>
    <td>必填</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>code</td>
    <td>int</td>
    <td>是</td>
    <td>返回码</td>
    <td>0</td>
  </tr>
  <tr>
    <td>msg</td>
    <td>string</td>
    <td>是</td>
    <td>错误信息</td>
    <td>null</td>
  </tr>
  <tr>
    <td>data</td>
    <td>object</td>
    <td>是</td>
    <td>正常信息</td>
    <td> {}</td>
  </tr>
</table>

**data 中的参数**
<table>
  <tr>
    <td>名称</td>
    <td>子属性</td>
    <td>类型</td>
    <td>描述</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>instance_list</td>
    <td></td>
    <td>[]</td>
    <td></td>
    <td></td>
  </tr>
  <tr>
    <td></td>
    <td>id</td>
    <td>string</td>
    <td>序号</td>
    <td>1</td>
  </tr>
    <td></td>
    <td>cluster_name</td>
    <td>string</td>
    <td>集群名</td>
    <td>gf.bridgx.online</td>
  </tr>
  </tr>
    <td></td>
    <td>instance_id</td>
    <td>string</td>
    <td>yyyy-dd-mm</td>
    <td>2021-10-11</td>
  </tr>
  </tr>
    <td></td>
    <td>startup_at</td>
    <td>string</td>
    <td>开机时间</td>
    <td>2021-11-11 15:15:20</td>
  </tr>
  </tr>
    <td></td>
    <td>shutdown_at</td>
    <td>string</td>
    <td>开机时间</td>
    <td>2021-11-11 15:45:20</td>
  </tr>
  </tr>
    <td></td>
    <td>startup_time</td>
    <td>int</td>
    <td>机器服务时长，单位：秒</td>
    <td>1800</td>
  </tr>
  </tr>
    <td></td>
    <td>instance_type</td>
    <td>string</td>
    <td>机器类型</td>
    <td>esc.c6.large</td>
  </tr>
  </tr>
    <td>pager</td>
    <td></td>
    <td>分页信息</td>
    <td></td>
    <td></td>
  </tr>
</table>

**请求示例**
<table>
  <tr>
    <td>名称</td>
    <td>示例值</td>
  </tr>
  <tr>
    <td>cluster_name</td>
    <td>gf.bridgx.online</td>
  </tr>
  <tr>
    <td>date</td>
    <td>2021-10-11</td>
  </tr>
  <tr>
    <td>page_number</td>
    <td>1</td>
  </tr>
  <tr>
    <td>page_size</td>
    <td>15</td>
  </tr>
</table>

**响应示例**

正常返回结果:
```JSON
{
    "code":200,
    "data":{
        "instance_list":[
            {
                "id":"1754",
                "cluster_name":"gf.scheduler.test",
                "instance_id":"i-2zeaavw***it89yvqm",
                "startup_at":"2021-11-18 10:16:53 +0800 CST",
                "shutdown_at":"2021-11-18 11:21:27 +0800 CST",
                "startup_time":3874,
                "instance_type":"ecs.s6-c1m1.small"
            },
            {
                "id":"1755",
                "cluster_name":"gf.scheduler.test",
                "instance_id":"i-2zeaavwojb***t89yvqn",
                "startup_at":"2021-11-18 10:16:53 +0800 CST",
                "shutdown_at":"2021-11-18 11:21:27 +0800 CST",
                "startup_time":3874,
                "instance_type":"ecs.s6-c1m1.small"
            },
            {
                "id":"1756",
                "cluster_name":"gf.bridgx.online",
                "instance_id":"i-2ze3ccvjn****1d9tzd0",
                "startup_at":"2021-11-18 10:31:53 +0800 CST",
                "shutdown_at":"2021-11-18 11:05:57 +0800 CST",
                "startup_time":2044,
                "instance_type":"ecs.s6-c1m1.small"
            }
        ],
        "pager":{
            "page_number":1,
            "page_size":10,
            "total":8
        }
    },
    "msg":"success"
}
```
异常返回结果：
```JSON
{
    "code":400,
    "msg":"param_invalid",
    "data":null
}
```
**返回码解释**

<table>
  <tr>
    <td>返回码</td>
    <td>状态</td>
    <td>解释</td>
  </tr>
  <tr>
    <td>200</td>
    <td>success</td>
    <td>执行成功</td>
  </tr>
  <tr>
    <td>400</td>
    <td>param_invalid</td>
    <td>参数有误</td>
  </tr>
</table> 



