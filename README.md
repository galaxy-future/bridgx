![z备份 12](https://user-images.githubusercontent.com/94337797/142638151-d38ff88d-e2ad-427d-bef5-2c0345557920.png)
======

[![Go Report Card](https://goreportcard.com/badge/github.com/galaxy-future/BridgX)](https://goreportcard.com/report/github.com/galaxy-future/BridgX) &nbsp;
[![CodeFactor](https://www.codefactor.io/repository/github/galaxy-future/bridgx/badge)](https://www.codefactor.io/repository/github/galaxy-future/bridgx)

BridgX是业界领先的基于全链路Serverless技术的云原生基础架构解决方案，目标是让开发者可以以开发单机应用系统的方式，开发跨云分布式应用系统，在不增加操作复杂度的情况下，兼得云计算的可扩展性和弹性伸缩等好处。将使企业用云成本减少50%以上，研发效率提升10倍。

它具有如下关键特性:

1、具备10分钟扩容1000台服务器的弹性能力；

2、一个平台统一管理不同云厂商的云上资源；

3、简洁易用，轻松上手；


联系我们
----

微信公众号：GalaxyFutureTech <br>
![image](https://user-images.githubusercontent.com/94337797/142592631-0bed59e6-7840-4c1c-870e-13dd9edd0c9b.png)

企业微信交流群:<br>
![image](https://user-images.githubusercontent.com/94337797/144558612-e7f36bd4-4afd-45ba-aa18-8be6d39c4537.png)




上手指南
----
#### 1、配置要求  
为了系统稳定运行，建议系统型号**2核4G内存**；BridgX已经在Linux系统以及macOS系统进行了安装和测试。


#### 2、环境依赖
- 如果已安装 Docker-1.10.0和Docker-Compose-1.6.0以上版本, 请跳过此步骤；如果没有安装，请查看[Docker Install](https://www.docker.com/products/container-runtime) 和 [Docker Compose Install](https://docs.docker.com/compose/install/);
- 如果已安装Git，请跳过此步骤；如果没有安装，请参照[Git - Downloads](https://git-scm.com/downloads)进行安装.


#### 3、安装部署  
下面是快速部署系统的步骤:

* (1)源码下载
  - 后端工程：
  > git clone https://github.com/galaxy-future/BridgX.git
  - 前端工程：
  > git clone https://github.com/galaxy-future/BridgX_FE.git

* (2)macOS系统部署
  - 后端部署,在BridgX目录下运行
    > make docker-run-mac
  - 前端部署,在BridgX_FE目录下运行
    > make docker-run-mac
   
  - 系统运行，在浏览器中输入 http://127.0.0.1 可以看到管理控制台界面,初始用户名root和密码为123456。

* (3)Linux安装部署
  - 以下步骤请使用 root用户 或有sudo权限的用户 sudo su - 切换到root用户后执行。
  - 1）针对使用者
    - 后端部署,在BridgX目录下运行,
      > make docker-run-linux
    - 前端部署,在BridgX_FE目录下运行
      > make docker-run-linux
    - 系统运行，浏览器输入 http://127.0.0.1 可以看到管理控制台界面,初始用户名root和密码为123456。

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

    - 前端部署
      - 如果跟后端同机部署，可以直接进入下一步;如果后端单独部署，则到 `cd conf` 下修改对应的配置ip和port信息。
      - 进入BridgX_FE目录下，然后使用以下命令
        > docker-compose up -d //启动BridgX前端服务 <br>
        > docker-compose down //停止BridgX前端服务  <br>

      - 系统运行，浏览器输入 http://127.0.0.1 可以看到管理控制台界面,初始用户名root和密码为123456。


    
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
