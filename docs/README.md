# Tiny-TikTok
## 1 项目介绍

### 1.1 概述

简易抖音项目后端实现

### 1.2 成员介绍

Ben

Marine

## 2 项目概览

### 2.1 技术选型与设计
[img](https://v1rwxew1bdp.feishu.cn/space/api/box/stream/download/asynccode/?code=MmJlMjRhMzJjZDc0OGNjYTIyY2JhOGQzMGQ2YmM1OGRfT2NKdmZWYkVnalpxbHBmeUtCTE1HUnB3djY2SEtUYW5fVG9rZW46U2kxbGJQOVE3b0hOWnB4M1p4bWNJZ2pKbjliXzE2OTE3NTA4MTg6MTY5MTc1NDQxOF9WNA)

* Gin：Web框架。高性能、轻量级、简洁，被广泛用于构建RESTful API、网站和其它HTTP服务

* **JWT：**用于身份的验证和授权，具有跨平台、无状态的优点，不需要再会话中保存任何信息，减轻服务器的负担

* **Hystrix：**服务熔断，防止由于单个服务的故障导致整个系统的崩溃

* **Etcd + grpc：**etcd实现服务注册和服务发现，grpc负责通信，构建健壮的分布式系统

* **OSS：**对象存储服务器，用于存储和管理非结构化数据，比如图片、视频等

* **FFmpeg：**多媒体开源工具，本项目用于上传视频时封面的截取

* **Gorm：**ORM库，用于操作数据库

* **MySQL：**关系型数据库，用于存储结构化数据

* **Redis：**键值存储数据库，以内存作为数据存储介质，提供高性能读写

* **RabbitMQ：**消息中间件，用于程序之间传递消息，支持消息的发布和订阅，支持消息异步传递，实现解耦和异步处理

  

### 2.2 总体设计

![img](https://v1rwxew1bdp.feishu.cn/space/api/box/stream/download/asynccode/?code=ZGE3ODFmZDU2YmNiM2VhMTdjZmQ1NWUxNDY1MzlkYWFfOFVQajI1WWJSeG4yalJ6UG1aSkpQYnhGVUdvVVk0TVdfVG9rZW46T1M0RmJFVFpob0lqUjR4cndvWmNKZTRPbjdiXzE2OTE3NTA4MzY6MTY5MTc1NDQzNl9WNA)

* 请求到达服务器前会对token进行校验
* 通过Api_Router对外暴露接口，进入服务，网关微服务对其它服务进行服务熔断和服务限流
* 各个服务先注册进入ETCD，api_router对各个服务进行调用，组装信息返回给服务端
* api_router通过gprc实现服务之间的通讯
* 服务操作数据库，将信息返回给上一层



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

​		**/cmd：**一个项目可以有很多个组件，吧main函数所在文件夹同一放在/cmd目录下

​		**/internal：**存放私有应用代码。handler类似三层架构中的控制层，service类似服务层，路由不用操作数据库，所以没有持久层

​		**/pkg：**存放可以被外部使用的代码库。auth中存放token鉴权，res存放对服务端的统一返回



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



* **工具函数包**

```bash
├─utils  
│  ├─etcd             # etcd服务注册与发现组件
│  ├─exception		 
│  └─snowFlake        # 雪花算法生成ID
```



### 2.4 数据库设计

![img](https://v1rwxew1bdp.feishu.cn/space/api/box/stream/download/asynccode/?code=NjQ1NDIzMDQ3YTUwY2U0ZDllNzA5ODc5YjEyMTIyZDNfa0NKQlVCMEZkaWpNMWRJYVdsWEFONlVJb0JaeExLdTFfVG9rZW46VXRMVGI3Qjg1b1VKUjN4T3JNRWN3bkxIbjNrXzE2OTE3NTA4NDY6MTY5MTc1NDQ0Nl9WNA)

**用户表：**用于存储用户名称、密码、头像、背景、用户简介信息，以由雪花算法生成的分布式id作为主键（其余表的ID同理），密码由bcrypt函数进行加密。

**视频表：**用于存储视频的作者、标题、封面路径、视频路径、获赞数量以及评论数量，以视频作者的id关联用户表。

**评论表：**用于存储视频的评论信息、评论创建时间、评论状态，通过用户id以及视频id关联用户表和视频表，通过评论状态作为软删除判断当前评论是否存在。

**消息表：**用于存放用户发送的消息以及消息的创建时间，通过用户id关联用户表，记录消息的发生者和消息的接收者

**关注表：**用户存放用户的关注信息，通过用户id关联用户表获取关注者和被关注者

**点赞表：**用于存放视频的点赞信息，通过用户id关联用户表获取点赞的用户，视频id关联视频表获取被点赞的视频



## 3 详细设计

## 4 测试

## 5 项目演示

## 6 总结





