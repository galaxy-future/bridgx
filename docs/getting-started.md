# 概览
本文是ComandX的快速入门指南，典型的用户包括四个步骤：<br>
1. 在云厂商账户模块，添加云厂商账户信息；<br>
2. 在集群管理模块，创建需要的集群模板；<br>
3. 在扩缩容任务模块，根据集群模板建立扩容任务；<br>
4. 在K8s集群模块，在扩容的机器上创建K8s集群及申请实例；<br>

# 前置条件
1、已经有云厂商账户，并且获得了AccessKey和AccessKey Secret信息，如果没有，请提前申请；<br>
[阿里云申请入口链接](https://help.aliyun.com/document_detail/53045.html)<br>

2、已经获取了网址信息，以及用户名和密码信息；<br>
如果是初次部署，可在浏览器中打开 http://127.0.0.1:80 登录系统,采用的认证方式是默认配置可用以下帐号登录<br>
默认用户名:root<br>
默认密码为:123456<br>

# 第一步：添加云厂商账户

1. 进入云厂商账户，点击添加云账户
![image](https://user-images.githubusercontent.com/94337797/142158688-a3a17da1-a068-4396-81fb-cf1f1270f184.png)

2. 填写云厂商的ak和sk信息，并保存；
![image](https://user-images.githubusercontent.com/94337797/142158808-19166f17-9ed6-4f5e-9ffe-65f698bbe7ed.png)


# 第二步：创建机型模板

1. 集群管理->创建集群
![image](https://user-images.githubusercontent.com/94337797/142158959-889069f7-1620-4b27-9764-2c6224b1ce72.png)

2. 进入云厂商配置页面，进行配置
![image](https://user-images.githubusercontent.com/94337797/145344282-048440e9-0005-40db-a39f-ecae0c6bd038.png)


3. 点击下一步进入网络配置页面：
![image](https://user-images.githubusercontent.com/94337797/145344359-e5a05478-0d6d-4473-afe2-b47830156903.png)


4. 网络配置完成后，点击下一步，进行机器规格配置
![image](https://user-images.githubusercontent.com/94337797/145344954-acb888f2-92c1-4060-899c-af3718c57bfe.png)


5. 机器规格配置完成后，点击下一步，进行存储磁盘配置，配置完成后，点击提交，创建集群成功
![image](https://user-images.githubusercontent.com/94337797/145345058-82aa08bd-ad3b-46d6-874e-1f25648d5895.png)


6. 点击完成按钮，则页面跳转到集群列表，有最新的创建的集群信息，则表示创建成功
![image](https://user-images.githubusercontent.com/94337797/145345871-d8e04f06-ed72-4f58-8505-dfb9ab5fa1e8.png)


# 第三步：创建快速扩缩容任务

1. 扩缩容任务->创建任务
![image](https://user-images.githubusercontent.com/94337797/145346132-372c2611-47dd-42e9-af1c-bd890bdc7319.png)


2. 填写扩缩容任务配置，填写完后点击提交按钮
![image](https://user-images.githubusercontent.com/94337797/145346567-69f33c5d-1aae-4693-b1a6-bebbd263fba1.png)


3. 提交成功后，页面跳转到任务列表，显示刚才创建的任务的执行情况
![image](https://user-images.githubusercontent.com/94337797/145346741-10b4494a-9a19-47bf-8200-ca911e1ddc54.png)


4. 如果需要查看本次任务创建的机器的IP，可点击执行明细查看详细信息
![image](https://user-images.githubusercontent.com/94337797/145346855-14bca836-0b33-43c2-99dd-9bbe07332df7.png)

# 第四步：创建K8s集群及申请实例

1. K8s集群->集群列表
![image](https://user-images.githubusercontent.com/94337797/145348389-667877c0-50e8-4201-bda0-e2c5ff3d77b0.png)

2. 点击创建集群，进入创建集群页面，填写配置信息
![image](https://user-images.githubusercontent.com/94337797/145348554-2e02790a-6614-4323-b431-cbc97c96e088.png)

3. 配置完成后，点击开始创建按钮，页面跳转到集群列表页面，显示集群的状态
![image](https://user-images.githubusercontent.com/94337797/145349578-227627ec-c289-4601-a6f3-72d54d2e8a78.png)

4. K8s集群创建成功后，在K8s集群->申请实例页面，申请实例
![image](https://user-images.githubusercontent.com/94337797/145353434-83190fef-cb4c-40d8-8426-871fb55f0201.png)



