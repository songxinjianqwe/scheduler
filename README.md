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
写写并发：submit任务之后的execute与stop不会在同一个goroutine中执行，由此会带来并发问题。<br />读写并发：List读取操作并不要求返回快照，否则代价太大，需要加全局锁，阻塞所有写操作（然后拷贝一份）。使用go test -race会检测到一系列的读写race，基本都是List()的读操作与对某些task的状态的写操作造成的，但这并不代表程序状态错误。<br />sync.Map#Range能做到的应该是遍历目前现存的所有元素，但不会保证每个元素的值都是在同一时刻的快照。而我们确实也不需要保证读到的是快照。
#### long polling（watch）
watch是自己实现了一个HTTP长轮询

## 待改进
* 支持多种script，如js、python、groovy
* 服务端对任务结果数保留一定数量，避免OOM
* 将shell命令放在沙箱环境中运行


