# BridgX
![z备份 15](https://user-images.githubusercontent.com/94337797/142117961-24f18cbe-4c00-4b57-86b2-8fe3c8fe2a92.png)

BridgX是业界领先的基于全链路Serverless技术的云原生基础架构解决方案，目标是让开发者可以以开发单机应用系统的方式，开发跨云分布式应用系统，在不增加操作复杂度的情况下，兼得云计算的可扩展性和弹性伸缩等好处。将使企业用云成本减少50%以上，研发效率提升10倍。

它具有如下关键特性:

1、具备10分钟扩容1000台服务器的弹性能力；

2、一个平台统一管理不同云厂商的云上资源；

3、简洁易用，轻松上手；

上手指南
----
#### 1、配置要求  
为了系统稳定运行，建议系统型号**2核4G内存及以上配置**。

#### 2、安装部署  
下面是快速部署系统的步骤:
* (1) 安装 **Docker-1.10.0**和**Docker-Compose-1.6.0**以上版本, 详细请查看[Docker Install](https://www.docker.com/products/container-runtime) 和 [Docker Compose Install](https://docs.docker.com/compose/install/),并确保**docker-daemon**已正常启动.(注意安装Docker Compose等组件的时候，由于国内网络环境的原因，访问github可能失败，遇此情况需要重试安装。)

* (2)确保**Git**已经安装，然后下载源码
  - 后端工程：
  > git clone https://github.com/galaxy-future/BridgX.git
  - 前端工程：
  > git clone https://github.com/galaxy-future/BridgX_FE.git
  - 由于项目会下载所需的必需基础镜像,建议将下载源码放到空间大于10G以上的目录中。 
* (3)以下步骤请使用 **root用户** 或有sudo权限的用户 **sudo su -** 切换到**root**用户后执行。
* (4)后端部署
  - 修改配置
    * BridgX依赖mysql和etcd组件，如果使用内置的mysql和etcd，可以跳过修改配置步骤。
    * 如果使用外部已搭建好的mysql和etcd，可以到 **cd conf** 下修改对应的配置信息。
  - 启动 BridgX，在根目录下，运行 
    >docker-compose up -d <br>
    * 如果已有外部mysql和etcd服务，并根据步骤修改配置做了变更，只需运行
    >docker-compose up -d api <br>
    >docker-compose up -d scheduler <br>
  - 停止 BridgX，在根目录下，运行
    >docker-compose down
* (5)前端部署
  * 修改配置
    * 如果跟后端同机部署，可以跳过修改配置步骤。<br>
    * 如果后端单独部署，可以到 **cd conf** 下修改对应的配置信息。
  * 启动 BridgX_FE，在根目录下，运行
    >docker-compose up -d <br>
    * BridgX_FE启动完成后，访问网址
http://localhost:8888 可以看到管理控制台界面,初始用户名root和密码为123456。
  * 停止 BridgX_FE，在根目录下，运行
    >docker-compose down


#### 3、快速上手  
通过[快速上手指南](https://github.com/galaxy-future/BridgX/blob/master/docs/getting-started.md)，可以掌握基本的快速扩缩容操作流程。  


#### 4、用户手册  
通过[用户手册](https://github.com/galaxy-future/BridgX/blob/master/docs/user-manual.md)，用户可以掌握BridgX的功能使用全貌，方便快速查找使用感兴趣的功能。

#### 5、开发者API手册
通过[开发者API手册](https://github.com/galaxy-future/BridgX/blob/master/docs/developer_api.md)，用户可以快速查看各项开发功能的API接口和调用方法，使开发者能够将BridgX集成到第三方平台上。

行为准则
------
[贡献者公约](https://github.com/galaxy-future/BridgX/blob/master/CODE_OF_CONDUCT)

授权
-----

BridgX使用[Apache License 2.0](https://github.com/galaxy-future/BridgX/blob/master/LICENSE)授权协议进行授权
