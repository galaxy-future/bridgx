![z备份 12](https://user-images.githubusercontent.com/94337797/142638151-d38ff88d-e2ad-427d-bef5-2c0345557920.png)
======

[![Go Report Card](https://goreportcard.com/badge/github.com/galaxy-future/BridgX)](https://goreportcard.com/report/github.com/galaxy-future/BridgX) &nbsp;
[![CodeFactor](https://www.codefactor.io/repository/github/galaxy-future/bridgx/badge)](https://www.codefactor.io/repository/github/galaxy-future/bridgx)


Language
----

English | [中文](https://github.com/galaxy-future/bridgx/blob/dev/docs/CH-README.md)

Introduction
-----
BridgX is the industry's leading cloud-native infrastructure solution based on full-process Serverless technology. The goal is to allow developers to create cross-cloud distributed application systems through building standalone application systems, without having to increase operational complexity, leveraging on both the scalability and elasticity of cloud computing. 
It has the following key features:
1. Has the flexibility to expand to 1000 servers within 1 minutes;
2. Supports K8s partitioning;
3. Provides a complete set of API ports;


Contact Us
----


[Weibo](https://weibo.com/galaxyfuture) | [Zhihu](https://www.zhihu.com/org/xing-yi-wei-lai) | [Bilibili](https://space.bilibili.com/2057006251)
| [WeChat Official Account](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat_official_account.md)
| [WeCom Communication Group](https://github.com/galaxy-future/comandx/blob/main/docs/resource/wechat.md)


Getting Started Guide
----
#### 1. Configuration Requirements

For stable operation of the system, the recommended system model is 2 CPU cores and 4G RAM; for Linux and macOS systems, BridgX has already been installed and tested.



#### 2. Environmental Dependence

- If you have already installed Docker-1.10.0 and Docker-Compose-1.6.0 or above, please skip this step; otherwise, please check 
[Docker Install](https://www.docker.com/products/container-runtime) 和 [Docker Compose Install](https://docs.docker.com/compose/install/);
- If you have already installed Git, please skip this step; otherwise, please refer to 
[Git - Downloads](https://git-scm.com/downloads) to install it.


#### 3. Installation and Deployment


* (1) Source code download
  - Back-end project：
  > git clone https://github.com/galaxy-future/bridgx.git
 
* (2) macOS system deployment
  - For back-end deployment,, execute in the BridgX directory:
    > make docker-run-mac
  - While the system is running, enter http://127.0.0.1 in the browser to see the management console interface. The initial user name is "root" and the password is "123456".

* (3) Linux installation and deployment
  - For the following steps, please use the "root" user or a user with "sudo" privileges, and run "sudo su-" to switch to the root user and execute it.

  - 1) For users 
    - For back-end deployment, execute in the BridgX directory:

      > make docker-run-linux
 
  - 2）For developers
    - Since the project will download the required basic image, we recommend placing the downloaded source code in a directory with more than 10G available storage space.

    - Back-end deployment

      - BridgX is dependent on the "mysql" and "etcd" components.
           - If you are using the built-in "mysql" and "etcd", then enter the BridgX root directory, and use the following command:
             > docker-compose up -d    //Start BridgX <br>
             > docker-compose down    //Stop BridgX  <br>
           - If you already have external "mysql" and "etcd" services, you can go to cd conf to modify the corresponding IP and port configuration information, and then use the following command while in the root directory of BridgX:
             > docker-compose up -d api    //Start api service <br>
             > docker-compose up -d scheduler //Start the scheduling service <br>
             > docker-compose down     //Stop BridgX service
#### 4.Developer's API Manual
Through the [Developers' API Manual](https://github.com/galaxy-future/bridgx/blob/dev/docs/en-developer-api.md), users can quickly view the API ports and calling methods of various developer functions, enabling developers to integrate BridgX into third-party platforms.


#### 5.Front-end Interface Operation
If you need to perform front-end operations, please install 
[ComandX](https://github.com/galaxy-future/comandx/blob/main/docs/EN-README.md).

Video Tutorial
------
[ComandX Installation](https://www.bilibili.com/video/BV1n34y167o8/) <br>
[Adding The Cloud Vendor Account](https://www.bilibili.com/video/BV1Jr4y1S7q4/)  <br>
[Create Cluster](https://www.bilibili.com/video/BV1Wb4y1v7jw/)   <br>
[Manual Scaling](https://www.bilibili.com/video/BV1bm4y197QD/)  <br>
[K8s Cluster Creation and Pod cutting](https://www.bilibili.com/video/BV1FY411p7rE/) <br>


Technical Articles
------
[《How does cloud-native technology migrate TB data per minute?》](https://zhuanlan.zhihu.com/p/442746588)<br>
[《Best practices for enterprise migration to K8s》](https://zhuanlan.zhihu.com/p/445131885) <br>
[《Top ten methods of cloud-native cost optimization 》](https://zhuanlan.zhihu.com/p/448405809)<br>

Code of Conduct
------
[Contributor Convention](https://github.com/galaxy-future/bridgx/blob/master/CODE_OF_CONDUCT)

Authorization
-----

BridgX uses [Apache License 2.0](https://github.com/galaxy-future/bridgx/blob/master/LICENSE) licensing agreement for authorization

Contact us
-----
If you want more information about the service, scan the following QR code to contact us:

![image](https://user-images.githubusercontent.com/102009012/163559389-813afa06-924f-412d-8642-1a0944384f91.png)


