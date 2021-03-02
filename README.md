# Mo2
[![codecov](https://codecov.io/gh/Monkey-Mouse/mo2/branch/main/graph/badge.svg?token=8X3HF5VFWT)](https://codecov.io/gh/Monkey-Mouse/mo2)
![gobadge](https://github.com/Monkey-Mouse/mo2/actions/workflows/.github/workflows/go.yml/badge.svg)
![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/Monkey-Mouse/mo2.svg)](https://github.com/Monkey-Mouse/mo2)
![GoReportCard](https://goreportcard.com/badge/github.com/Monkey-Mouse/mo2)
[![BCH compliance](https://bettercodehub.com/edge/badge/Monkey-Mouse/mo2?branch=main)](https://bettercodehub.com/)  
一个属于所有人的博客网站
![logo](mo2front/public/img/icons/android-chrome-512x512.png)

- [Mo2](#mo2)
  - [Why](#why)
  - [项目结构](#项目结构)
  - [运行项目](#运行项目)
    - [使用docker](#使用docker)
  - [从源码编译](#从源码编译)
    - [先决条件](#先决条件)
    - [编译服务器](#编译服务器)
    - [编译前端](#编译前端)

## Why
在现在的网络上，有众多的公开博客网站，大家可以在上边发表自己的文章。  

但是这些网站都更新的很慢，风格千篇一律，很多用户想要的功能无法及时添加，用户无法获取自己真正最想要的体验。  

同时，大部分博客网站用户是希望能建造自己的个人博客的，或者说用户希望使用一个自己建造的博客。然而大部分人没有精力这么做，即使有word press这样的框架帮助，能自建博客的人任然很少，何况word press建造的基础网站本身也存在很多问题。  

所以我们的项目就是为了解决这个问题而诞生：我们意图打造一个开源的博客系统，使用持续集成和持续部署技术，吸引有想自己设计自己的博客的用户为我们的仓库提交代码。在用户的代码被接受并合并到生产分支后，我们的持续部署技术会在10分钟以内更新我们的服务器，让用户看到自己的更改，同时我们会记录下贡献者的信息，对他们进行鸣谢。让所有参与网站开发的人都有成就感。  

所以其实这个项目只有一小半是博客，一大半是社区。我们的目标是靠这个项目的特性吸引高质量用户，形成社区，将MO2建设为  
> A blog site made for everyone, and made by everyone  

## 项目结构
如你所见，最外层是一个go项目，我们的前端项目在[mo2front](/mo2front/)目录下。我们的计划是让每个重要的文件夹下
都有个readme文件。  
## 运行项目
运行项目有两种方法，第一种是使用docker，第二种是从源码进行编译。  

### 使用docker
如果你只是想运行或部署Mo2项目而不对他的源码感兴趣，使用docker无疑是你最好的选择。  
- 首先，你需要安装[docker](https://docs.docker.com/engine/install/)以及[docker-compose](https://docs.docker.com/compose/install/)  
- 然后，你需要下载我们[updateWatcher](updateWatcher)文件夹下的所有文件，维持它们的相对位置
- 最后，在此目录打开命令行，执行命令`docker-compose up`即可，此时将可以通过http://localhost:5001/swagger/index.html 访问后端控制台

> **注意** 这样运行后端有部分功能不能使用。包括：
> - image upload相关功能
> - email发送相关功能
> 
> 使用这些功能需要预先设定特殊的环境变量


## 从源码编译
### 先决条件
- go 1.15
- 尽可能新的 npm 包管理工具
- 尽可能新的 mongodb

### 编译服务器
首先，我们需要确保你的mongodb已经在运行。  
然后，我们需要把你的mongodb地址导出到环境变量中，
如果您的mongodb是默认配置的话，你可以使用以下命令：  
linux bash:  
```bash
export MO2_MONGO_URL=mongodb://127.0.0.1:27017
```
windows powershell:
```powershell
$env:MO2_MONGO_URL=mongodb://127.0.0.1:27017
```
然后，使用同一个终端，在项目根目录：
```bash
go run main.go
```
即可运行后端
> **注意** 这样运行后端有部分功能不能使用。包括：
> - image upload相关功能
> - email发送相关功能
> 
> 使用这些功能需要预先设定特殊的环境变量

### 编译前端
在项目根目录打开一个终端，然后切换到[mo2front](/mo2front/)目录  
```bash
cd ./mo2front/
```
然后，使用npm进行install：
```bash
npm install
```
最后，运行前端
```bash
npm run serve
```
根据命令行提示打开对应的网页就可以看到前端页面
> **注意** Mo2在默认配置下，前端的运作并不要求本地运行后端，我们的设置会让npm在本地进行代理，将请求代理到我们的生产环境的服务器上。如果你想更改这项设置，请更改[vue.config.js](mo2front/vue.config.js)





