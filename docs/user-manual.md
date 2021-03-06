# BridgX用户手册
## 1. 引言
### 1.1 标识
  本文档是BridgX开源云原生系统的用户使用说明文档，对应初始开源版本V1.0。  
       
### 1.2 系统概述
  BridgX是业界领先的基于全链路Serverless技术的云原生基础架构解决方案，目标是让开发者可以以开发单机应用系统的方式，开发跨云分布式应用系统，在不增加操作复杂度的情况下，兼得云计算的可扩展性和弹性伸缩等好处。将使企业用云成本减少50%以上，研发效率提升10倍。  
  
  BridgX是团队结合此前在微博建设10年左右的后端体系经验，与业界最新的技术如Serverless、混合云、K8s、AIOps、Service Mesh等进行整合，通过在业务代码中嵌入BridgX SDK，可为业务提供RPC、缓存、存储等各类函数服务，并通过BridgX Agent代理向BridgX Server发起请求，BridgX Server可实现多云对接、容器编排、AIOps等高级功能。
  
  客户通过使用BridgX，可实现全链路托管，不用再操心基础架构，实现单机式开发，通过函数访问所有服务，并且可实现全链路弹性，具备实时评估十分钟扩容万台的能力，全自动运维，避免人为失误。除此之外，客户还可以灵活选择云厂商与云服务商，数据可灵活迁移，保证数据资产安全可控，通过智能调度，使得客户用云成本下降50%以上。
### 1.3 体系架构
  BridgX系统由多个模块构成，用户可以通过管理系统统一管理。Dashboard模块是用户方便的查看使用的云资源情况，包括集群、任务、机器等信息；集群管理模块可以使用户根据业务需要定义适配的集群模板，并且可以对集群模板进行管理；扩缩容模块可以使用户在有扩容或缩容需求时，快速建立扩缩容任务，并且可以对扩容容任务进行管理；云厂商账户模块可以帮助客户管理不同云厂商的账户信息；云服务器模块帮助客户管理自己申请的服务器资源；费用管理模块使用户可以看到自已云资源的成本费用情况；账户管理模块可以帮助用户管理子账户，对子账户进行增删改查。
### 1.4 文档概述
  本文是BridgX系统的用户手册，详细介绍了BridgX系统的部署和使用方法。欢迎任何人阅读和使用本文档，并将其提供给其他对BridgX平台感兴趣的同学。
### 1.5 名词解释

<table>
  <tr>
    <td>序号</td>
    <td>术语</td>
    <td>定义</td>
  </tr>
  <tr>
    <td>1</td>
    <td>集群</td>
    <td>一组配置相同机器的共同模板</td>
  </tr>
  <tr>
    <td>2</td>
    <td>扩容</td>
    <td>增加云服务器的台数</td>
  </tr>
  <tr>
    <td>3</td>
    <td>缩容</td>
    <td>减少云服务器的台数</td>
  </tr>
  <tr>
    <td>4</td>
    <td>任务</td>
    <td>一个扩容或缩容的作业任务</td>
  </tr>
       
</table>




## 2. 软件综述
### 2.1 软件应用
  BridgX系统可用于构建组织内部IT基础设施，提供快速扩缩容、集群管理、任务管理、机器管理、费用管理等服务。BridgX系统可以利用用户自己拥有的物理服务器构建的私有云，也可以使用阿里云公有云平台实现上述功能。
  BridgX系统提供了丰富的配置和二次开发API接口，允许用户根据自己的实际应用情况对系统进行调整，使BridgX系统可以支持广泛的用户需求。
### 2.2 用户角色说明
  管理员、开发人员、业务人员
### 2.3 授权和使用
  BridgX遵循Apache License 2.0开源协议。用户可以自由下载并免费使用此系统。用户需要自行承担此使用此系统带来的风险。
## 3. 扩缩容任务
  扩缩容任务模块主要用于帮助用户根据自己的业务需求快速进行定量的扩容或缩容任务，并对任务的执行情况进行管理。
### 3.1 创建任务
  创建任务模块是进行扩缩容时的便捷入口，当客户有扩缩容任务时，选择好扩缩容的集群模板、扩缩容的数量以及类型，即可快速进行扩缩容任务。当创建的任务没有合适的集群时，也可以通过添加集群的快速入口进行配置。
![image](https://user-images.githubusercontent.com/94337797/142165917-9b4db33f-1544-4546-9dc3-9197d83f1084.png)


### 3.2 任务列表
  任务列表模块可以帮助用户对已经建立的任务进行管理，可以根据任务名、集群、状态进行搜索和筛查，快速查找感兴趣的任务。同时，通过新建按钮，也可以快速进入创建任务页面，进行新的任务创建。
![image](https://user-images.githubusercontent.com/94337797/142166055-d73d9f06-d2da-4496-8e26-459a0b4f4e98.png)


  如果对任务的描述感兴趣，可以通过点击任务名查看，包括任务的ID、集群模块、变更动作、变更的机器数核执行时间等信息。
![image](https://user-images.githubusercontent.com/94337797/142166097-398baafe-a3c2-4864-aca9-7846a16f2736.png)



  如果对任务的执行细节感兴趣，可以点击执行明细，进行便捷查看，包括机器的IP、状态、开始时间、耗时等信息。
  ![image](https://user-images.githubusercontent.com/94337797/142166153-7ff0eed7-3a14-4308-9d4c-c6e1fd6b0a0e.png)


## 4. 集群管理
  集群管理模块主要用于创建配置相同机器的配置模板，包括云厂商、网络配置、存储配置、系统镜像、机器规格等相同参数，并对创建的集群模板进行管理，增删查改等。
### 4.1 创建集群
  创建集群主要包括云厂商配置、网络配置、机器规格和磁盘4个步骤。  
  
  云厂商配置主要是进行集群名称、集群描述的填写以及云厂商和云账户的选择。再进行云账户选择时，应提前在云厂商账户模块提前录入用户的云账户信息。
  ![image](https://user-images.githubusercontent.com/94337797/142166250-8d59762e-6873-468e-b37e-eb6920bbdcbe.png)



  网络配置模块主要是对集群所属的云厂商的可用区域、可用区进行选择，以及所在的vpc网络、子网配置、安全组配置。如果还有合适的VPC、子网或安全组，则可以通过相应的左边的按钮快捷创建。如果需要进行公网访问，可以可以开启公网配置，设置带宽收费类型和最大的带宽。
   ![image](https://user-images.githubusercontent.com/94337797/142166287-a8ed9cd5-9193-4c62-b00e-2ae3f20f571e.png)
 


  机器规格配置包括机型配置和镜像配置，机型配置是指选择需要的计算核数和内存大小，镜像配置是选择需要的操作系统镜像。
 ![image](https://user-images.githubusercontent.com/94337797/142166339-631eb4ea-7e25-48aa-b99b-55031ba5591d.png)

  磁盘配置包括对系统盘和数据盘的配置，可以对系统盘的类型和容量进行配置，对数据盘则可以配置类型、容量以及数据盘的个数。
 ![image](https://user-images.githubusercontent.com/94337797/142166686-b3b5b151-a6a6-464f-bc76-2ea3bf966876.png)

  
### 4.2 集群列表
  集群列表提供对已经创建的集群进行管理的功能，可以通过集群名称、云厂商或Access Key信息进行选择和筛查，而且通过创建集群按钮，可以快速的进入创建集群页面，添加新的集群。
  ![image](https://user-images.githubusercontent.com/94337797/142167306-766b91b2-55e9-4ac1-90ad-c68e1b96e6a9.png)


  如果对某个集训需要进行更改或删除，则可以先在列表的第一列进行选择，然后点击相应的编辑或删除按钮。
  ![image](https://user-images.githubusercontent.com/94337797/142167370-6e882ca6-b844-46bf-9c87-56266bee6e90.png)


 
## 5. 云服务器
  云服务器模块，可以显示用户所有已经申请的云服务器信息，并且可以通过机器名、IP、云厂商等信息进行选择和筛查。
  ![image](https://user-images.githubusercontent.com/94337797/142167437-44eb6e9c-366f-4368-8d64-03aeccbd080e.png)


  如果对某个特定的机器详细感兴趣，可以通过点击相应的机器名称进行详细信息查看，包括实例ID、创建时间、云厂商、机器规格、镜像ID、系统盘类型及大小、数据盘类型大小及个数，网络配置包括VPC名称、子网名称、安全组名称、内网IP及公网IP等信息。
  ![image](https://user-images.githubusercontent.com/94337797/142167482-4855aed0-b930-426d-800a-648a61e7ab1a.png)



## 6. 云厂商账户
  云厂商账户模块主要是管理客户的云账户信息，可以通过账户名称、云厂商、账户信息等进行查找和筛选。
  ![image](https://user-images.githubusercontent.com/94337797/142167534-2bda3f7c-9cee-49ce-abd1-e220f6c6f563.png)


  如果需要增加云账户，则可以点击添加云账户按钮，填写账户名称、云厂商以及账户的AccessKey和AccessKey Secret信息。
  ![image](https://user-images.githubusercontent.com/94337797/142167574-cc0a21b2-9472-4116-9f02-ffa6dd0d8c45.png)


  如果需要对某个账户进行删除或者修改，则可以通过勾选相应的账户，然后对点击删除或者编辑按钮，进行相应的操作。
  ![image](https://user-images.githubusercontent.com/94337797/142167633-0790eb96-a59e-4b70-87b0-0bfa71b88464.png)



## 7. 费用管理
  费用模块主要是用来对客户已经申请的云资源所花费的费用进行管理。客户可以查看所有机器的费用，也可以通过集群选择查看某个集群的费用，同时支持可以针对某天的费用进行查询。
  ![image](https://user-images.githubusercontent.com/94337797/142167675-8f167fc7-7ff1-44ae-bba1-1c98d4ff6af7.png)



## 8. 账户管理
  账户管理模块主要是用于对主账户和子账户进行进行管理，可以创建子账户、对子账户进行禁用和启用。
  ![image](https://user-images.githubusercontent.com/94337797/142167718-6bcdf7eb-9be0-40f9-8c34-76a526f31c6b.png)


  当需要创建子账号时，点击创建子账号按钮，通过添加用户名、密码来添加子账号。
  ![image](https://user-images.githubusercontent.com/94337797/142167769-9d6c9677-c97c-468c-8290-2d29ce320938.png)

