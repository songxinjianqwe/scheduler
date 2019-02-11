# Scheduler

# 单机scheduler实现，包含客户端和服务器
[https://github.com/songxinjianqwe/scheduler](https://github.com/songxinjianqwe/scheduler)
## 初衷
1、毕设计划实现一个简化版的Docker容器<br />2、目前需要熟悉Go语言<br />3、学习Docker的架构
* Docker分为Docker Daemon（HTTP Server）和Docker CLI（命令行工具），而我自己实现的scheduler也是这种架构

4、后续可以将命令执行环境从宿主机移至容器中，以增强隔离性。计划依赖containerd来实现，进一步熟悉容器技术。
## 架构
CS架构，Client为CLI工具，Server为HTTP Server，均使用Go语言编写。

## 安装
### 客户端CLI

1. `go get "github.com/songxinjianqwe/scheduler"`
1. `cd $GOPATH/src/github.com/songxinjianqwe/scheduler/cli`
1. `go install`
1. .`/cli`

### 服务器

1. `cd $GOPATH/src/github.com/songxinjianqwe/scheduler/daemon`
1. `go install`
1. `./daemon`<br />
启动服务器后会在命令行中打印出REST API：

* 返回所有任务，对应客户端list命令<br />
ROUTE: /api/tasks<br />
Path regexp: ^/api/tasks[/]?$<br />
Queries templates:<br />
Queries regexps:<br />
Methods: GET
* 返回该id对应的任务，对应客户端get命令<br />
ROUTE: /api/tasks/{id}<br />
Path regexp: /]+)[/]?$<br />
Queries templates:<br />
Queries regexps:<br />
Methods: GET
* 提交任务，对应客户端submit命令<br />
ROUTE: /api/tasks<br />
Path regexp: ^/api/tasks[/]?$<br />
Queries templates:<br />
Queries regexps:<br />
Methods: POST
* 停止任务，对应客户端stop命令<br />
ROUTE: /api/tasks/{id}<br />
Path regexp: /]+)[/]?$<br />
Queries templates:<br />
Queries regexps:<br />
Methods: PUT
* 删除任务，对应客户端delete命令<br />
ROUTE: /api/tasks/{id}<br />
Path regexp: /]+)[/]?$<br />
Queries templates:<br />
Queries regexps:<br />
Methods: DELETE

## 功能

### 查询任务列表
#### `命令`
`./cli list`
#### 示例![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549696338600-7738aeee-e986-406c-bf1b-e01aad22ce8c.png#align=left&display=inline&height=62&linkTarget=_blank&name=image.png&originHeight=78&originWidth=1259&size=24090&width=1007)<br />
### 提交任务
#### 命令
`./cli submit task_id -type=task_type -time=delay_or_interval -script="shell script"`
* task_id：必填，任务ID，要求是全局唯一
* -type：必填，任务类型；可以是delay，表示延迟任务，也可以是cron，表示定时任务
* -time：必填，延时时间；格式类似于10s,1m,20h等，可以参考Go的time.Duration字符串格式
* -script：必填，shell脚本，如果脚本中含有空格，则需要用引号包裹
#### 示例
* `./cli submit print_ls_per_10s  -type=cron -time=10s -script="ls"`
  * 表示提交一个任务ID为print_ls_per_10s的定时任务，每隔10秒，会在服务器上执行一次ls命令
* `./cli submit do_calc_after_1min  -type=delay -time=1m -script="sleep 10s;echo $((1+2))"`
  * 表示提交一个任务ID为do_calc_after_1min的延时任务，在1分钟后，会睡眠10秒，然后打印出1+2的结果
### 读取/监听单个任务
#### 命令
`./cli get task_id [-watch=true]`
* task_id：任务ID，要求是全局唯一
* -watch：选填，是否需要监听；默认为false，如果设置为true，那么get会不断返回该任务的最新状态
#### 示例
* 。`./cli get do_calc_after_1mins`

![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549697164813-a392c7df-71ff-4ec0-bd7c-b45e4c56b747.png#align=left&display=inline&height=206&linkTarget=_blank&name=image.png&originHeight=257&originWidth=622&size=47560&width=498)
* `./cli get do_calc_after_1mins -watch=true`

![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549697218319-10eae443-4aa8-4365-8857-3ad62f64b45c.png#align=left&display=inline&height=579&linkTarget=_blank&name=image.png&originHeight=724&originWidth=733&size=117884&width=586)
### 停止任务
#### 命令
`./cli stop task_id`<br />停止任务：
* 如果是延迟任务，则：
  * 可以在任务开始前停止，则任务不会被执行，终态为Stopped
  * 如果任务已经开始执行，则无法停止，且报错，终态为Executed
  * 如果任务已经执行完毕，则无法停止，且报错，终态为Executed
* 如果是定时任务，则：
  * 可以在任务开始前停止，则任务一次都不会被执行，终态为Stopped
  * 如果在某一次任务开始执行后停止，则任务本次执行不会被中止，且会保存本次执行结果，终态为Stopped
  * 如果在某一个任务任务执行后，下一次任务执行前停止，则下次执行会被跳过，终态为Stopped
#### 示例
![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549697344116-e5db240a-305f-4fe4-b56e-001eb26b32e5.png#align=left&display=inline&height=85&linkTarget=_blank&name=image.png&originHeight=106&originWidth=554&size=27500&width=443)![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549697390057-7a94bb20-58d1-4b0f-972e-a27e2ce103ba.png#align=left&display=inline&height=47&linkTarget=_blank&name=image.png&originHeight=62&originWidth=582&size=12730&width=442)
### 删除任务
#### 命令
`./cli delete task_id`
#### 示例
![image.png](https://cdn.nlark.com/yuque/0/2019/png/257642/1549697473486-b2eb5e58-486d-4883-ada4-23cd33254402.png#align=left&display=inline&height=84&linkTarget=_blank&name=image.png&originHeight=105&originWidth=602&size=19755&width=482)
## 难点
#### 并发
写写并发：submit任务之后的execute与stop不会在同一个goroutine中执行，由此会带来并发问题。<br />~~读写并发：List读取操作并不要求返回快照，否则代价太大，需要加全局锁，阻塞所有写操作（然后拷贝一份）。使用go test -race会检测到一系列的读写race，基本都是List()的读操作与对某些task的状态的写操作造成的，但这并不代表程序状态错误。~~<br />~~sync.Map#Range能做到的应该是遍历目前现存的所有元素，但不会保证每个元素的值都是在同一时刻的快照。而我们确实也不需要保证读到的是快照。~~
##### 首次检测race

```text

==================
WARNING: DATA RACE
Read at 0x00c0001b62b8 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:86 +0x89
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001b62b8 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:174 +0x27b
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0001b62d0 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:86 +0x89
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001b62d0 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).updateStatus()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:190 +0x42
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:177 +0x2f2
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0001b62d8 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:86 +0x89
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001b62d8 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).updateStatus()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:191 +0x8d
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:177 +0x2f2
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0001b62f0 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:86 +0x89
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001b62f0 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).increaseVersionAndSignalListeners()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:196 +0x63
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:185 +0x303
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0001a0780 by goroutine 15:
  reflect.typedmemmove()
      /usr/local/Cellar/go/1.11.5/libexec/src/runtime/mbarrier.go:177 +0x0
  reflect.packEface()
      /usr/local/Cellar/go/1.11.5/libexec/src/reflect/value.go:119 +0x103
  reflect.valueInterface()
      /usr/local/Cellar/go/1.11.5/libexec/src/reflect/value.go:1008 +0x16f
  reflect.Value.Interface()
      /usr/local/Cellar/go/1.11.5/libexec/src/reflect/value.go:978 +0x51
  encoding/json.marshalerEncoder()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:448 +0x99
  encoding/json.(*structEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:647 +0x307
  encoding/json.(*structEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:661 +0x7b
  encoding/json.(*arrayEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:769 +0x12a
  encoding/json.(*arrayEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:776 +0x7b
  encoding/json.(*sliceEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:743 +0xf1
  encoding/json.(*sliceEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:755 +0x7b
  encoding/json.(*structEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:647 +0x307
  encoding/json.(*structEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:661 +0x7b
  encoding/json.(*arrayEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:769 +0x12a
  encoding/json.(*arrayEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:776 +0x7b
  encoding/json.(*sliceEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:743 +0xf1
  encoding/json.(*sliceEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:755 +0x7b
  encoding/json.(*encodeState).reflectValue()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:333 +0x93
  encoding/json.(*encodeState).marshal()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:305 +0xad
  encoding/json.Marshal()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:160 +0x73
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:30 +0xfd
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001a0780 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:174 +0x22e
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0001a0798 by goroutine 15:
  reflect.Value.String()
      /usr/local/Cellar/go/1.11.5/libexec/src/reflect/value.go:1711 +0x5c
  encoding/json.stringEncoder()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:610 +0xda
  encoding/json.(*structEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:647 +0x307
  encoding/json.(*structEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:661 +0x7b
  encoding/json.(*arrayEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:769 +0x12a
  encoding/json.(*arrayEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:776 +0x7b
  encoding/json.(*sliceEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:743 +0xf1
  encoding/json.(*sliceEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:755 +0x7b
  encoding/json.(*structEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:647 +0x307
  encoding/json.(*structEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:661 +0x7b
  encoding/json.(*arrayEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:769 +0x12a
  encoding/json.(*arrayEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:776 +0x7b
  encoding/json.(*sliceEncoder).encode()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:743 +0xf1
  encoding/json.(*sliceEncoder).encode-fm()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:755 +0x7b
  encoding/json.(*encodeState).reflectValue()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:333 +0x93
  encoding/json.(*encodeState).marshal()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:305 +0xad
  encoding/json.Marshal()
      /usr/local/Cellar/go/1.11.5/libexec/src/encoding/json/encode.go:160 +0x73
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:30 +0xfd
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0001a0798 by goroutine 23:
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:174 +0x22e
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 23 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
--- FAIL: TestSubmit (3.01s)
    testing.go:771: race detected during execution of test
time="2019-02-11T09:42:54+08:00" level=info msg="Receive a task:common.Task{Id:\"test_delay_7ebfd222-c94c-4c71-b923-88122c260b54\", TaskType:\"delay\", Time:2000000000, Script:\"echo 1\", Results:[]common.TaskResult(nil), Status:0, LastStatusUpdated:time.Time{wall:0xc6dbb70, ext:63685446174, loc:(*time.Location)(0x17d8900)}, Version:0, timer:(*time.Timer)(nil), ticker:(*time.Ticker)(nil), lock:(*sync.Mutex)(nil), watchCond:(*sync.Cond)(nil)}"
time="2019-02-11T09:42:56+08:00" level=info msg="start executing task[test_delay_7ebfd222-c94c-4c71-b923-88122c260b54]"
time="2019-02-11T09:42:56+08:00" level=info msg="Task stdout: 1\n"
time="2019-02-11T09:42:57+08:00" level=info msg="start stopping task[test_delay_7ebfd222-c94c-4c71-b923-88122c260b54]"
FAIL
coverage: 60.3% of statements
Found 6 data race(s)

```


仔细观察，发现其中很大一部分是先写（比如task的状态更新），再在List()中json序列化时读，出现读写冲突。<br />但是在List()时遍历理论上是拷贝了一份的！<br />为了我编写了一个示例：

```go
type Person struct {
	Name string
	Age int
	Cars []Car
}

type Car struct {
	Name string
}

func createPerson(Name string, Age int, carNames []string) *Person {
	person := Person{}
	person.Name = Name
	person.Age = Age
	var cars []Car
	for _, carName := range carNames {
		cars = append(cars, createCar(carName))
	}
	person.Cars = cars
	return &person
}

func createCar(name string) Car {
	car := Car{}
	car.Name = name
	return car
}
var allPersons = []*Person{
	createPerson("p1", 1, []string{"c1","c2"}),
	createPerson("p2", 2, []string{"c1","c2"}),
}

func getPersonList() []Person {
	var personList []Person
	for _, p := range allPersons {
		personList = append(personList, *p)
	}
	return personList
}

func TestCopy(t *testing.T) {
	list := getPersonList()
	for _, p := range list {
		fmt.Printf("%#v\n", p)
	}
	allPersons[0].Name="p3"
	allPersons[0].Cars[0].Name = "c3"
	for _, p := range list {
		fmt.Printf("%#v\n", p)
	}
}
```
输出结果：
> main.Person{Name:"p1", Age:1, Cars:[]main.Car{main.Car{Name:"c1"}, main.Car{Name:"c2"}}}
> main.Person{Name:"p2", Age:2, Cars:[]main.Car{main.Car{Name:"c1"}, main.Car{Name:"c2"}}}
> main.Person{Name:"p1", Age:1, Cars:[]main.Car{main.Car{Name:"c3"}, main.Car{Name:"c2"}}}
> main.Person{Name:"p2", Age:2, Cars:[]main.Car{main.Car{Name:"c1"}, main.Car{Name:"c2"}}}


注意，这里修改person的Name，并没有影响到拷贝后的personList；但是修改car的状态，会影响！<br />原因是[]Car是一个切片类型，切片类型是引用类型，而struct结构体拷贝是浅拷贝，所以没有拷贝[]Car！
##### 再次检测race
这是List()原来的实现，我
```go
// List返回的是原来的一份拷贝
func (this *StandAloneEngine) List() ([]common.Task, error) {
	var tasks []common.Task
	this.tasks.Range(func(key, value interface{}) bool {
		tasks = append(tasks, *value.(*common.Task))
		return true
	})
	if tasks == nil {
		tasks = make([]common.Task, 0)
	}
	return tasks, nil
}


```
然后我修改了第5行：<br />`tasks = append(tasks, value.(*common.Task).Clone())`

```go
func (this *Task) Clone() Task {
	aCopy := *this
	aCopy.lock = nil
	aCopy.ticker = nil
	aCopy.timer = nil
	aCopy.watchCond = nil
	aCopy.Results = make([]TaskResult, len(this.Results))
	copy(aCopy.Results, this.Results)
	return aCopy
}
```
再次检测race，又发现了问题：

```text

==================
WARNING: DATA RACE
Read at 0x00c0000c0b78 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:201 +0xe2
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0000c0b78 by goroutine 22:
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:174 +0x27b
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 22 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0000c0b90 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:201 +0xe2
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0000c0b90 by goroutine 22:
  github.com/songxinjianqwe/scheduler/common.(*Task).updateStatus()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:190 +0x42
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:177 +0x2f2
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 22 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0000c0b98 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:201 +0xe2
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0000c0b98 by goroutine 22:
  github.com/songxinjianqwe/scheduler/common.(*Task).updateStatus()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:191 +0x8d
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:177 +0x2f2
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 22 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0000c0bb0 by goroutine 15:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:201 +0xe2
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0000c0bb0 by goroutine 22:
  github.com/songxinjianqwe/scheduler/common.(*Task).increaseVersionAndSignalListeners()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:196 +0x63
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:185 +0x303
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 22 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
==================
WARNING: DATA RACE
Read at 0x00c0000892c0 by goroutine 15:
  runtime.slicecopy()
      /usr/local/Cellar/go/1.11.5/libexec/src/runtime/slice.go:221 +0x0
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:207 +0x211
  sync.(*Map).Range()
      /usr/local/Cellar/go/1.11.5/libexec/src/sync/map.go:337 +0x13c
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).List()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:85 +0x77
  github.com/songxinjianqwe/scheduler/daemon/handler.GetAllTasksHandler()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:25 +0xa1
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a

Previous write at 0x00c0000892c0 by goroutine 22:
  github.com/songxinjianqwe/scheduler/common.(*Task).appendResultAndUpdateStatusAtomically()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:174 +0x22e
  github.com/songxinjianqwe/scheduler/common.(*Task).Execute()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/common/task.go:97 +0x2a6
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit.func1()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:57 +0x62

Goroutine 15 (running) created at:
  net/http.(*Server).Serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2851 +0x4c5
  net/http.(*Server).ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2764 +0xe8
  net/http.ListenAndServe()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:3004 +0xef
  github.com/songxinjianqwe/scheduler/daemon/server.Run()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/server/server.go:49 +0x522

Goroutine 22 (finished) created at:
  github.com/songxinjianqwe/scheduler/daemon/engine/standalone.(*StandAloneEngine).Submit()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/engine/standalone/standalone_engine_impl.go:55 +0x3c0
  github.com/songxinjianqwe/scheduler/daemon/handler.SubmitTask()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/daemon/handler/handler.go:85 +0x35f
  net/http.HandlerFunc.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1964 +0x51
  github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux.(*Router).ServeHTTP()
      /Users/jasper/go/src/github.com/songxinjianqwe/scheduler/vendor/github.com/gorilla/mux/mux.go:212 +0x12e
  net/http.serverHandler.ServeHTTP()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:2741 +0xc4
  net/http.(*conn).serve()
      /usr/local/Cellar/go/1.11.5/libexec/src/net/http/server.go:1847 +0x80a
==================
--- FAIL: TestSubmit (3.01s)
    testing.go:771: race detected during execution of test
time="2019-02-11T16:34:33+08:00" level=info msg="Receive a task:common.Task{Id:\"test_delay_266aa0c3-fc83-4fa7-abb6-1ed69d692ae3\", TaskType:\"delay\", Time:2000000000, Script:\"echo 1\", Results:[]common.TaskResult(nil), Status:0, LastStatusUpdated:time.Time{wall:0xf9dbde0, ext:63685470873, loc:(*time.Location)(0x17d7640)}, Version:0, timer:(*time.Timer)(nil), ticker:(*time.Ticker)(nil), lock:(*sync.RWMutex)(nil), watchCond:(*sync.Cond)(nil)}"
time="2019-02-11T16:34:35+08:00" level=info msg="start executing task[test_delay_266aa0c3-fc83-4fa7-abb6-1ed69d692ae3]"
time="2019-02-11T16:34:35+08:00" level=info msg="Task stdout: 1\n"
time="2019-02-11T16:34:36+08:00" level=info msg="start stopping task[test_delay_266aa0c3-fc83-4fa7-abb6-1ed69d692ae3]"
FAIL
Found 5 data race(s)
FAIL	github.com/songxinjianqwe/scheduler/cli/client	9.045s
?   	github.com/songxinjianqwe/scheduler/cli/command	[no test files]
?   	github.com/songxinjianqwe/scheduler/common	[no test files]
?   	github.com/songxinjianqwe/scheduler/common/run	[no test files]
?   	github.com/songxinjianqwe/scheduler/daemon	[no test files]
?   	github.com/songxinjianqwe/scheduler/daemon/engine	[no test files]
ok  	github.com/songxinjianqwe/scheduler/daemon/engine/standalone	43.060s
?   	github.com/songxinjianqwe/scheduler/daemon/handler	[no test files]
?   	github.com/songxinjianqwe/scheduler/daemon/server	[no test files]
```


其中有sliceCopy的读与task的appendResult之间的读写冲突，于是我在Clone外面加了一层读锁，以此保证读到的results是在一次完整的写操作之后、下次写操作之前的一个正确的快照。
```go
func (this *Task) Clone() Task {
	this.lock.RLock()
	defer this.lock.RUnlock()
	aCopy := *this
	aCopy.lock = nil
	aCopy.ticker = nil
	aCopy.timer = nil
	aCopy.watchCond = nil
	aCopy.Results = make([]TaskResult, len(this.Results))
	copy(aCopy.Results, this.Results)
	return aCopy
}
```
此时再次检测race就没有问题了！

#### long polling（watch）
watch是自己实现了一个HTTP长轮询

## 待改进
* 支持多种script，如js、python、groovy
* 服务端对任务结果数保留一定数量，避免OOM
* 将shell命令放在沙箱环境中运行


