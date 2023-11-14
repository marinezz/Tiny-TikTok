<h1 align = "center">Tiny-TikTok</h1>

<p align="center">
    <a href="https://img.shields.io/badge/lan-go-green">
        <img alt="Static Badge" src="https://img.shields.io/badge/lan-go-green">
    </a>
    <a href="https://img.shields.io/badge/web-gin-blue">
       <img alt="Static Badge" src="https://img.shields.io/badge/web-gin-blue">
    </a>
	<a href="https://img.shields.io/badge/orm-gorm-yellow">
       <img alt="Static Badge" src="https://img.shields.io/badge/orm-gorm-yellow">
    </a>
	<a href="https://img.shields.io/badge/database-MySQL-red">
       <img alt="Static Badge" src="https://img.shields.io/badge/database-MySQL-red">
    </a>
	<a href="https://img.shields.io/badge/cache-Redis-pink">
       <img alt="Static Badge" src="https://img.shields.io/badge/cache-Redis-pink">
    </a>
</p>

## 1 项目介绍

简易抖音项目后端实现，使用 **Gin** 作为web框架，**MySQL** 作为数据存储并使用 Gorm 操作数据库。整个项目分为用户服务、视频服务、社交服务，使用 **Etcd** 作为注册中心，**Grpc** 进行服务之间的通信。 采用 **Redis** 作为缓存，提高读写效率；使用消息中间件 **RabbitMQ**，达到上游服务和下游服务的解耦。

![image](https://github.com/marinezz/Tiny-TikTok/blob/main/docs/image/%E9%A1%B9%E7%9B%AE%E8%8E%B7%E5%A5%96%E8%AF%81%E4%B9%A6.png?raw=true)

**简单写一个总结：**

青训营结束也过去了一段时间，这次获得了**三等奖**，还算比较满意的结果（自己离还有不少差距，心服口服）。本来两个人完成的这个项目，自己也不是主力golang选手，一路跌跌撞撞，也将近用了完整的三四周时间。

收获最大的有两个点。一是只有思维的摩擦，才会有创造，比如缓存redis的设计，一开始只想放在后端和数据库之间，简单降低数据库的压力，仅此而已。但是后面在思考的时候，面对不同的场景，采用不同的数据结构，以及如何保证数据的一致性等等。第二个就是实践，只有真正去写代码了，才能意识到很多问题，理论上能滔滔不绝的说出来，但是实践出来还是有很大的差距，一些细节根本在想的时候抓不住。

最后总结一下自己的优势与不足，也算给后续想参加青训的同志们的一些简易吧。本项目对比其余项目，可能（自我感觉）优于其它人的，就是测试比较完整。对于redis、MQ的使用，这些大家基本都会去使用，但是使用的时候一定要往深度去思考，比如redis中数据的一致性（这个问题在最后答辩的时候都会问，有的人会考虑掉或者做得不够好）。然后就是广度，这是我们明显缺失的，对于前面的我最佩服的几个项目，从技术选型，到团队合作（一些团队合作工作），项目部署（自动化部署），服务监控，日志等都几乎完整的实现。而我们部署靠手动，监控只统一处理了日志，还有很大很大的差距。还有技术选型有一点点小问题，因为我参加过上一届青训，但是那次没项目就纯看视频课（没啥用），所以这次的课没有咋看，这次是教了Hertz和kitex的，所以可能会给出一种没好好听课的感觉吧~

**如果对你有帮助，右上角 Satr 走起 ！！！**

### 1.1 目录介绍

**第一部分：项目简介**

**第二部分：项目概览。包括整个项目的技术选择、整体设计、目录结构设计与数据库设计**

**第三部分：项目详细设计。包括对项目的思考以及一些关键点的设计思路**

**第四部分：项目测试**

**第五部分：启动说明。拿到整个项目，应该如何让它跑起来**

<br>

### 1.2 分支介绍

本项目的分支是项目的几个迭代版本，可以根据新老版本查看，由简到难：

* **main分支**：最新版本，仍在迭代中。修复一些小bug，加入消息队列RabbitMQ

* **1.0分支**：最老的分支，项目完整的启动，最简单的版本

* **2.0分支**：在1.0的基础上优化缓存，索引。

<br>

## 2 项目概览

### 2.1 技术选型与设计

![image](https://github.com/marinezz/Tiny-TikTok/blob/main/docs/image/技术架构图.png)

* **Gin**：Web框架。高性能、轻量级、简洁，被广泛用于构建RESTful API、网站和其它HTTP服务

* **JWT**：用于身份的验证和授权，具有跨平台、无状态的优点，不需要再会话中保存任何信息，减轻服务器的负担

* **Hystrix**：服务熔断，防止由于单个服务的故障导致整个系统的崩溃

* **Etcd + grpc**：etcd实现服务注册和服务发现，grpc负责通信，构建健壮的分布式系统

* **OSS**：对象存储服务器，用于存储和管理非结构化数据，比如图片、视频等

* **FFmpeg**：多媒体开源工具，本项目用于上传视频时封面的截取

* **Gorm**：ORM库，用于操作数据库

* **MySQL**：关系型数据库，用于存储结构化数据

* **Redis**：键值存储数据库，以内存作为数据存储介质，提供高性能读写

* **RabbitMQ**：消息中间件，用于程序之间传递消息，支持消息的发布和订阅，支持消息异步传递，实现解耦和异步处理

<br>

### 2.2 总体设计

![image](https://github.com/marinezz/Tiny-TikTok/blob/main/docs/image/%E6%80%BB%E4%BD%93%E8%AE%BE%E8%AE%A1%E5%9B%BE.png)

* 请求到达服务器前会对token进行校验
* 通过Api_Router对外暴露接口，进入服务，网关微服务对其它服务进行服务熔断和服务限流
* 各个服务先注册进入ETCD，api_router对各个服务进行调用，组装信息返回给服务端
* api_router通过gprc实现服务之间的通讯
* 服务操作数据库，将信息返回给上一层

<br>

### 2.3 项目结构设计

* **总目录结构**

```bash
├─api_router       # 路由网关
├─docs             # 项目文档
├─social_service   # 社交服务
├─user_service     # 用户服务
├─utils            # 工具函数包
└─video_service    # 视频服务
```

<br>

* **路由网关**

```bash
├─api_router   
│  ├─cmd  
│  ├─config       # 项目配置文件
│  ├─discovery     # 服务注册与发现
│  ├─internal
│  │  ├─handler
│  │  └─service
│  │      └─pb   
│  ├─pkg  
│  │  ├─auth
│  │  └─res
│  └─router        # 路由和中间件  
│      └─middleware
```

​		**/cmd**：一个项目可以有很多个组件，吧main函数所在文件夹同一放在/cmd目录下

​		**/internal**：存放私有应用代码。handler类似三层架构中的控制层，service类似服务层，路由不用操作数据库，所以没有持久层

​		**/pkg**：存放可以被外部使用的代码库。auth中存放token鉴权，res存放对服务端的统一返回

<br>

* **具体服务**

```bash
├─user_service
│  ├─cmd
│  ├─config
│  ├─discovery
│  ├─internal
│  │  ├─handler
│  │  ├─model
│  │  └─service       # 持久化层  
│  │      └─pb
│  └─pkg
│      └─encryption   # 密码加密
```

<br>

* **工具函数包**

```bash
├─utils  
│  ├─etcd             # etcd服务注册与发现组件
│  ├─exception		 
│  └─snowFlake        # 雪花算法生成ID
```

<br>

### 2.4 数据库设计

![image](https://github.com/marinezz/Tiny-TikTok/blob/main/docs/image/数据库设计图.png)

**用户表**：用于存储用户名称、密码、头像、背景、用户简介信息，以由雪花算法生成的分布式id作为主键（其余表的ID同理），密码由bcrypt函数进行加密。

**视频表**：用于存储视频的作者、标题、封面路径、视频路径、获赞数量以及评论数量，以视频作者的id关联用户表。

**评论表**：用于存储视频的评论信息、评论创建时间、评论状态，通过用户id以及视频id关联用户表和视频表，通过评论状态作为软删除判断当前评论是否存在。

**消息表**：用于存放用户发送的消息以及消息的创建时间，通过用户id关联用户表，记录消息的发生者和消息的接收者

**关注表**：用户存放用户的关注信息，通过用户id关联用户表获取关注者和被关注者

**点赞表**：用于存放视频的点赞信息，通过用户id关联用户表获取点赞的用户，视频id关联视频表获取被点赞的视频

<br>

## 3 详细设计

本部分包含：**认证鉴权**、**分布式唯一ID**、**密码加密**、**数据库操作**、**视频上传**、**日志打印**、**服务熔断**、**统一错误处理**、**高性能读写**、**异步解耦**的设计

详细设计参考答辩文档第三部分：[说明文档](https://v1rwxew1bdp.feishu.cn/docx/ATJPdobcOouDDLxVHsycpANMnig?from=from_copylink)

<br>

## 4 测试

本部分包含：单元测试、接口测试（功能测试）、性能测试（压力测试）

测试详情查看测试文档：[测试报告](https://v1rwxew1bdp.feishu.cn/docx/F1aQdSY6AoIzeLx6B8tcVQhgnud?from=from_copylink)

<br>

## 5 启动说明

### 5.1 启动之前

项目启动之前，了解各个组件的版本，一直觉得版本是一个很大的问题，先搞清楚每一个组件版本，后续才能有条不紊的进行

|     名称     | 版本  |         作用         |
| :----------: | :---: | :------------------: |
|    **Go**    | 1.19  |                      |
|  **MySQL**   |  8.0  |     关系型数据库     |
|   **Etcd**   | 3.5.1 |       注册中心       |
|  **Redis**   | 5.0.7 |         缓存         |
| **RabbitMQ** | 3.7.4 | 消息中间件，异步解耦 |
|  **FFmpeg**  |  6.0  |    视频剪切为封面    |
|  **Protoc**  | 3.15  |      生成pb文件      |

确保每一个组件安装成功，并将ffmpeg和protoc配置进入环境变量

 

### 5.2 拉取项目

保证安装了git，直接克隆到本地

```bash
git clone https://github.com/marinezz/Tiny-TikTok.git
```

<br>

### 5.3 拉取依赖

将所有需要的依赖拉取到本地，进入每一个文件，拉取依赖。以api_router为例（其余文件操作相同）：

```bash
cd api_router/
go mod tidy
```

<br>

### 5.4 启动项目

**第一步**：启动etcd。确保etcd注册中心启动成功。如果不放心可以下载ectdkeeper（参考网络），后续也可以参看服务是否注册进入注册中心

**第二步**：建立数据库。建立名为tiny_tiktok的数据库，不用建表，启动时gorm会自动建表

**第三步**：修改配置文件。参考配置文件example，根据自己的实际情况修改配置文件

**第四步**：正式启动项目。到每个文件的cmd文件中启动项目，或者使用命令 `go run main.go`  

**第五步**：项目启动成功

<br>

## 6 总结与展望

项目的从0到有，经过了一个月的时间，在开发过程中遇到了很多问题，也解决很多问题，得到了很多收获。在不断的实际操作中，才能有新的收获，新的想法。截止青训营的结束，现在还是有很多想法的存在，我们也会把自己的想法付诸实现，完完整整做出自己满意的项目。



我们还会继续ing.....



**如果对你有帮助的话，希望您不要吝啬你得star哦！！**

![image](https://github.com/marinezz/Tiny-TikTok/blob/main/docs/image/%E8%8E%B7%E5%A5%96%E8%AF%81%E4%B9%A6.png?raw=true)
