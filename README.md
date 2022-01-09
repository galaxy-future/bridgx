![z备份 12](https://user-images.githubusercontent.com/94337797/142638151-d38ff88d-e2ad-427d-bef5-2c0345557920.png)
======

[![Go Report Card](https://goreportcard.com/badge/github.com/galaxy-future/BridgX)](https://goreportcard.com/report/github.com/galaxy-future/BridgX) &nbsp;
[![CodeFactor](https://www.codefactor.io/repository/github/galaxy-future/bridgx/badge)](https://www.codefactor.io/repository/github/galaxy-future/bridgx)

语言(language)
----

[English](https://github.com/galaxy-future/bridgx/blob/dev/docs/EN-README.md)

简介
--------
BridgX是业界领先的基于全链路Serverless技术的云原生基础架构解决方案，目标是让开发者可以以开发单机应用系统的方式，开发跨云分布式应用系统，在不增加操作复杂度的情况下，兼得云计算的可扩展性和弹性伸缩等好处。

它具有如下关键特性:

1、具备1分钟扩容1000台服务器的弹性能力；

2、支持K8s切割；

3、提供完善的API接口；


联系我们
----
[微博](https://weibo.com/galaxyfuture) | [知乎](https://www.zhihu.com/org/xing-yi-wei-lai) | [B站](https://space.bilibili.com/2057006251)
| [微信公众号](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat_official_account.md)
| [企业微信交流群](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat.md)


上手指南
----
#### 1、配置要求  
为了系统稳定运行，建议系统型号**2核4G内存**；BridgX已经在Linux系统以及macOS系统进行了安装和测试。


#### 2、环境依赖
- 如果已安装 Docker-1.10.0和Docker-Compose-1.6.0以上版本, 请跳过此步骤；如果没有安装，请查看[Docker Install](https://www.docker.com/products/container-runtime) 和 [Docker Compose Install](https://docs.docker.com/compose/install/);
- 如果已安装Git，请跳过此步骤；如果没有安装，请参照[Git - Downloads](https://git-scm.com/downloads)进行安装.


#### 3、安装部署  

* (1)源码下载
  - 后端工程：
  > git clone https://github.com/galaxy-future/bridgx.git
 
* (2)macOS系统部署
  - 后端部署,在BridgX目录下运行
    > make docker-run-mac
  - 系统运行，在浏览器中输入 http://127.0.0.1 可以看到管理控制台界面,初始用户名root和密码为123456。

* (3)Linux安装部署
  - 以下步骤请使用 root用户 或有sudo权限的用户 sudo su - 切换到root用户后执行。
  - 1）针对使用者
    - 后端部署,在BridgX目录下运行,
      > make docker-run-linux
 
  - 2）针对开发者
    - 由于项目会下载所需的必需基础镜像,建议将下载源码放到空间大于10G以上的目录中。
    - 后端部署
      - BridgX依赖mysql和etcd组件，
           - 如果使用内置的mysql和etcd，则进入BridgX根目录，则使用以下命令：            
             > docker-compose up -d    //启动BridgX <br>
             > docker-compose down    //停止BridgX  <br>
           - 如果已经有了外部的mysql和etcd服务，则可以到 `cd conf` 下修改对应的ip和port配置信息,然后进入BridgX的根目录，使用以下命令:
             > docker-compose up -d api    //启动api服务 <br>
             > docker-compose up -d scheduler //启动调度服务 <br>
             > docker-compose down     //停止BridgX服务
#### 4、开发者API手册
通过[开发者API手册](https://github.com/galaxy-future/bridgx/blob/master/docs/developer_api.md)，用户可以快速查看各项开发功能的API接口和调用方法，使开发者能够将BridgX集成到第三方平台上。

#### 5、前端界面操作
如果需要进行前端操作，请安装[ComandX](https://github.com/galaxy-future/comandx/blob/main/README.md)

视频教程
------
[BridgX安装](https://www.bilibili.com/video/BV1n34y167o8/) <br>
[添加云账户](https://www.bilibili.com/video/BV1Jr4y1S7q4/)  <br>
[创建集群](https://www.bilibili.com/video/BV1Wb4y1v7jw/)   <br>
[手动扩缩容](https://www.bilibili.com/video/BV1bm4y197QD/)  <br>
[K8s集群创建与Pod切割](https://www.bilibili.com/video/BV1FY411p7rE/)<br>


技术文章
------
[《云原生技术如何每分钟级迁移TB级数据》](https://zhuanlan.zhihu.com/p/442746588)<br>
[《企业迁移到K8s的最佳实践》](https://zhuanlan.zhihu.com/p/445131885) <br>
[《来自一线大厂的十大云原生成本优化手段》](https://zhuanlan.zhihu.com/p/448405809)<br>



行为准则
------
[贡献者公约](https://github.com/galaxy-future/bridgx/blob/master/CODE_OF_CONDUCT)

授权
-----

BridgX使用[Apache License 2.0](https://github.com/galaxy-future/bridgx/blob/master/LICENSE)授权协议进行授权
