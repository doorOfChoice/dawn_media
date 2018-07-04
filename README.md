# dawndevil media house

这是属于dawndveil的视频网站，也是该死的课程设计。

写的仓促，BUG挺多....


## 功能

#### 管理员后台
* 用户增删改查
* 媒体文件增删改查
* 分类增删改查
* 评论查改
* 主页

#### 用户后台
* 用户修改密码
* 用户修改个人信息
* 浏览记录
* 点赞媒体
* 评论管理

#### 前台
* 主页
* 媒体多条件查找
* 影片观看，多清晰度切换
* 影片点赞
* 影片评论


## 配置文件

```toml
#端口地址
address=":8080"

#数据库账户
username="root"

#密码
password=""

#数据库
database=""

#注册密码salt
pass_salt="dawndevil"

#登录salt
login_salt="dawndevil"

#服务器生成的session对应cookie名称
session_name="media_web"

#头像目录
avatar_dir="static/media/avatar/"

#头像网址
avatar_map="/static/media/avatar/"

#封面目录
cover_dir="static/media/cover/"

#封面网址
cover_map="/static/media/cover/"

#媒体目录
media_dir="static/media/video/"

#媒体网址
media_map="/static/media/video/"

```

## 使用

```go
    #在GOPATH下Clone
    git clone github.com/doorOfChoice/dawn_media
    cd github.com/doorOfChoice/dawn_media
    go get .
    //======windows
    go build -o server.exe
    server.exe

    //======linux
    go build -o server
    ./server
```

## 页面
|页面|URI|
|:-:|:-:|
|管理员界面|localhost:8080/admin|
|用户界面|localhost:8080/ordinary|
|首页|localhost:8080|